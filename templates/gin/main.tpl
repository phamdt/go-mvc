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

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // blank import necessary to use driver
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"go.uber.org/zap"
)

func main() {
	// construct dependencies
	log := zap.NewExample().Sugar()
	defer log.Sync()

	// setup database
	db, err := newDb()
	if err != nil {
		log.Fatalf("can't initalize database connection: %v", zap.Error(err))
		return
	}

	// setup router and middleware
	router := controllers.GetRouter(log, db)
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// setup monitoring only if the license key is set
	nrKey := os.Getenv("NR_LICENSE_KEY")
	if nrKey != "" {
		nrMiddleware, err := newRelic(nrKey)
		if err != nil {
			log.Fatal("Unexpected error setting up new relic", zap.Error(err))
			panic(err)
		}
		router.Use(nrMiddleware)
	}

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
