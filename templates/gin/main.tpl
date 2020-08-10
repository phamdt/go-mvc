package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"{{Name}}/controllers"

	"github.com/getsentry/sentry-go"
  sentrygin "github.com/getsentry/sentry-go/gin"
  "github.com/gin-contrib/pprof"
  "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // blank import necessary to use driver
	"github.com/prometheus/client_golang/prometheus/promhttp"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"go.uber.org/zap"
)

func main() {
	// construct dependencies

	// setup app logging
	log := zap.NewExample().Sugar()
	defer log.Sync()

	// setup request logging separately
	requestLogger, _ := zap.NewProduction()

	// setup database
	db, err := newDb()
	if err != nil {
		log.Fatalf("can't initialize database connection: %v", zap.Error(err))
		return
	}

	// setup router and middleware
	router := controllers.GetRouter(log, db)

	// setup Sentry for monitoring
	if err := sentry.Init(sentry.ClientOptions{
	    Dsn: "your-public-dsn",
	}); err != nil {
	    log.Infof("Sentry initialization failed: %v\n", err)
	}
	sentryOptions := sentrygin.Options{
		// Whether Sentry should repanic after recovery, in most cases it should be set to true,
		// as gin.Default includes its own Recovery middleware that handles http responses.
		Repanic:					true,
		// Whether you want to block the request before moving forward with the response.
		// Because Gin's default `Recovery` handler doesn't restart the application,
		// it's safe to either skip this option or set it to `false`.
		WaitForDelivery: 	false,
		// Timeout for the event delivery requests.
		Timeout:         5 * time.Second,
	}
	router.Use(sentrygin.New(sentryOptions))

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	router.Use(ginzap.Ginzap(requestLogger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	router.Use(ginzap.RecoveryWithZap(requestLogger, true))

	// setup New Relic monitoring only if the license key is set
	nrKey := os.Getenv("NR_LICENSE_KEY")
	if nrKey != "" {
		nrMiddleware, err := newRelic(nrKey)
		if err != nil {
			log.Fatal("Unexpected error setting up new relic", zap.Error(err))
			panic(err)
		}
		router.Use(nrMiddleware)
	}

	// setup pprof and prometheus server separate from application server so as to
	// keep profiling information available only on localhost and not exposed to
	// the internet in production
	go func() {
		internalRouter := gin.Default()
		pprof.Register(internalRouter)

		internalRouter.Get("/metrics", gin.WrapH(promhttp.Handler))
		internalRouter.Run(":8081")
	}()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", zap.Error(err))
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Info("timeout of 5 seconds.")
	}
	log.Info("Server exiting")
}

func newRelic(nrKey string) (gin.HandlerFunc, error) {
	cfg := newrelic.NewConfig(os.Getenv("APP_NAME"), nrKey)
	// Creates a New Relic Application
	apm, err := newrelic.NewApplication(cfg)
	if err != nil {
		return nil, err
	}
	return nrgin.Middleware(apm), nil
}

func newDb() (*sqlx.DB, error) {
	configString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))
	return sqlx.Open("postgres", configString)
}
