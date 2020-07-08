package gomvc

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file3 := &embedded.EmbeddedFile{
		Filename:    "build/circleciconfig.yml.tpl",
		FileModTime: time.Unix(1587572809, 0),

		Content: string("version: 2\njobs:\n  build_and_test:\n    docker:\n      - image: circleci/golang:1.13\n    working_directory: /go/src/{{gitRepoPath}}\n    steps:\n      - checkout\n      - setup_remote_docker:\n          docker_layer_caching: true\n      - add_ssh_keys\n{{#envFileName}}\n      - run:\n          name: Add environment variables to a file\n          command: cp {{#envFileSampleName}} {{envFileName}}\n{{/envFileName}}\n      - run:\n          name: Start Containers\n          command: docker-compose -f docker-compose.yml up -d\n      - run:\n          name: Wait for Server\n          command: |\n            chmod +x .circleci/wait-for-server-start.sh\n            .circleci/wait-for-server-start.sh\n      - run:\n          name: Wait extra 10s to ensure database is seeded\n          command: sleep 10\n      - run:\n          name: Run tests\n          command: docker exec -it {{containerName}} go test ./...\n\nworkflows:\n  version: 2\n  build:\n    jobs:\n      - build_and_test"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "build/docker-compose.yml.tpl",
		FileModTime: time.Unix(1587572809, 0),

		Content: string("version: \"3\"\nservices:\n  {{Name}}_postgres:\n    container_name: {{Name}}_db\n    hostname: {{Name}}_db\n    image: \"postgres:11\"\n    env_file: .env\n    ports:\n      - \"5432:5432\"\n# UNCOMMENT ONCE YOU HAVE MIGRATIONS\n#  {{Name}}_migrations:\n#    container_name: migrations\n#    image: migrate/migrate:v4.6.2\n#    command: [\"-path\", \"/migrations/\", \"-database\", \"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable\", \"up\"]\n#    depends_on:\n#      - postgres\n#    env_file: .env\n#    restart: on-failure\n#    links: \n#      - postgres\n#    volumes:\n#      - ./migrations:/migrations \n#\n  {{Name}}:\n    container_name: {{Name}}\n    build:\n      context: .\n      dockerfile: Dockerfile\n    env_file: .env\n    volumes:\n      - ./:/go/src/{{Name}}\n    ports:\n      - \"8080:8080\"\n    links:\n      - {{Name}}_postgres\n\n  {{Name}}_test:\n    container_name: {{Name}}_test\n    build:\n      context: .\n      dockerfile: test.Dockerfile\n    env_file: .env\n    volumes:\n      - ./:/go/src/{{Name}}\n    ports:\n      - \"9999:9999\"\n    links:\n      - {{Name}}_postgres\n\n"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "build/env.tpl",
		FileModTime: time.Unix(1587572809, 0),

		Content: string("# Postgres Database\n# Env vars originate from the postgres image on dockerhub\nPOSTGRES_HOST={{Name}}\nPOSTGRES_USER={{Name}}\nPOSTGRES_DB={{Name}}\nPOSTGRES_PASSWORD={{Name}}\n\nAPP_NAME={{Name}}\nNR_LICENSE_KEY="),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "build/sqlboiler.toml.tpl",
		FileModTime: time.Unix(1587572809, 0),

		Content: string("[psql]\n  dbname = \"{{dbName}}\"\n  host   = \"0.0.0.0\"\n  port   = {{dbPort}}\n  user   = \"{{dbUser}}\"\n  pass   = \"{{dbPassword}}\"\n  blacklist = [\n    {{#blacklist}}{{blacklist}}{{/blacklist}}\n  ]\n  sslmode = \"disable\"\n{{#templates}}\n  templates = [\n    {{templates}}\n  ]\n{{/templates}}\n"),
	}
	file7 := &embedded.EmbeddedFile{
		Filename:    "build/wait-for-server-start.sh.tpl",
		FileModTime: time.Unix(1587572809, 0),

		Content: string("#!/bin/bash\n\necho \"Waiting for servers to start...\"\nattempts=1\nwhile true; do\n  docker exec -i {{Name}} curl -f http://localhost:8080/health > /dev/null 2> /dev/null\n  if [ $? = 0 ]; then\n    echo \"Service started\"\n    break\n  fi\n  ((attempts++))\n  if [[ $attempts == 5 ]]; then\n    echo \"5 attempts to check health failed\"\n    break\n  fi\n  sleep 10\n  echo $attempts\ndone"),
	}
	file9 := &embedded.EmbeddedFile{
		Filename:    "gin/controller.tmpl",
		FileModTime: time.Unix(1594235845, 0),

		Content: string("package controllers\n\nimport (\n\t\"net/http\"\n\n\t\"github.com/gin-gonic/gin\"\n\t\"github.com/jmoiron/sqlx\"\n\t{{# ORM }}\n\t\"github.com/volatiletech/sqlboiler/boil\"\n\t{{/ ORM }}\n\t\"go.uber.org/zap\"\n\t\"{{ModuleName}}/models\"\n)\n\n// {{Name}}Controller exposes the methods for interacting with the\n// RESTful {{Name}} resource\ntype {{Name}}Controller struct {\n\tdb  *sqlx.DB\n\tlog *zap.SugaredLogger\n}\n\n{{#each Actions}}\n{{{ whichAction Handler }}}\n{{/each}}\n\n{{#each ErrorResponses}}\nfunc (ctrl *{{../Name}}Controller) is{{Name}}(c *gin.Context) bool {\n\t// TODO: Add your controller-specific logic for determining if the request \n\t// should return a {{Name}} response with status code {{Code}} as\n\t// found in your spec: {{Ref}}\n\treturn false\n}\n{{/each}}\n"),
	}
	filea := &embedded.EmbeddedFile{
		Filename:    "gin/main.tpl",
		FileModTime: time.Unix(1594235941, 0),

		Content: string("package main\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"log\"\n\t\"net/http\"\n\t\"os\"\n\t\"os/signal\"\n\t\"syscall\"\n\t\"time\"\n\t\"{{Name}}/controllers\"\n\n\t\"github.com/gin-gonic/gin\"\n\t\"github.com/jmoiron/sqlx\"\n\t_ \"github.com/lib/pq\" // blank import necessary to use driver\n\tnewrelic \"github.com/newrelic/go-agent\"\n\t\"github.com/newrelic/go-agent/_integrations/nrgin/v1\"\n\t\"go.uber.org/zap\"\n)\n\nfunc main() {\n\t// construct dependencies\n\tlog := zap.NewExample().Sugar()\n\tdefer log.Sync()\n\n\t// setup database\n\tdb, err := newDb()\n\tif err != nil {\n\t\tlog.Fatalf(\"can't initialize database connection: %v\", zap.Error(err))\n\t\treturn\n\t}\n\n\t// setup router and middleware\n\trouter := controllers.GetRouter(log, sqlxDB)\n\t// Recovery middleware recovers from any panics and writes a 500 if there was one.\n\trouter.Use(gin.Recovery())\n\n\t// setup monitoring only if the license key is set\n\tnrKey := os.Getenv(\"NR_LICENSE_KEY\")\n\tif nrKey != \"\" {\n\t\tnrMiddleware, err := newRelic(nrKey)\n\t\tif err != nil {\n\t\t\tlog.Fatal(\"Unexpected error setting up new relic\", zap.Error(err))\n\t\t\tpanic(err)\n\t\t}\n\t\trouter.Use(nrMiddleware)\n\t}\n\n\tsrv := &http.Server{\n\t\tAddr:    \":8080\",\n\t\tHandler: router,\n\t}\n\n\tgo func() {\n\t\t// service connections\n\t\tif err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {\n\t\t\tlog.Fatalf(\"listen: %s\\n\", zap.Error(err))\n\t\t}\n\t}()\n\n\t// Wait for interrupt signal to gracefully shutdown the server with\n\t// a timeout of 5 seconds.\n\tquit := make(chan os.Signal)\n\t// kill (no param) default send syscall.SIGTERM\n\t// kill -2 is syscall.SIGINT\n\t// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it\n\tsignal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)\n\t<-quit\n\tlog.Info(\"Shutdown Server ...\")\n\n\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n\tdefer cancel()\n\tif err := srv.Shutdown(ctx); err != nil {\n\t\tlog.Fatal(\"Server Shutdown:\", zap.Error(err))\n\t}\n\t// catching ctx.Done(). timeout of 5 seconds.\n\tselect {\n\tcase <-ctx.Done():\n\t\tlog.Info(\"timeout of 5 seconds.\")\n\t}\n\tlog.Info(\"Server exiting\")\n}\n\nfunc newRelic(nrKey string) (gin.HandlerFunc, error) {\n\tcfg := newrelic.NewConfig(os.Getenv(\"APP_NAME\"), nrKey)\n\t// Creates a New Relic Application\n\tapm, err := newrelic.NewApplication(cfg)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn nrgin.Middleware(apm), nil\n}\n\nfunc newDb() (*sqlx.DB, error) {\n\tconfigString := fmt.Sprintf(\"host=%s user=%s dbname=%s password=%s\", os.Getenv(\"POSTGRES_HOST\"), os.Getenv(\"POSTGRES_USER\"), os.Getenv(\"POSTGRES_DB\"), os.Getenv(\"POSTGRES_PASSWORD\"))\n\treturn sqlx.Open(\"postgres\", configString)\n}\n"),
	}
	filec := &embedded.EmbeddedFile{
		Filename:    "gin/partials/create.tmpl",
		FileModTime: time.Unix(1594234969, 0),

		Content: string("// Create saves a new {{Resource}} record into the database\nfunc (ctrl *{{Resource}}Controller) {{Name}}(c *gin.Context) {\n\tm := models.{{SingularResource}}{}\n\tif err := c.ShouldBindJSON(m); err != nil {\n\t\tctrl.log.Error(\"invalid {{Resource}} creation request\",\n\t\t\tzap.Error(err),\n\t\t)\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t\treturn\n\t}\n\t{{# ORM }}\n\terr := m.Insert(ctrl.db, boil.Infer())\n\tif err != nil {\n\t\tctrl.log.Error(\"error creating {{Name}}\",\n\t\t\tzap.Error(err))\n\t\tc.AbortWithStatus(http.StatusInternalServerError)\n\t}\n\t{{/ ORM }}\n\tc.JSON(http.StatusCreated, gin.H{})\n}\n"),
	}
	filed := &embedded.EmbeddedFile{
		Filename:    "gin/partials/delete.tmpl",
		FileModTime: time.Unix(1594234969, 0),

		Content: string("// Delete deletes a new {{Resource}} record into the database\nfunc (ctrl *{{Resource}}Controller) {{Name}}(c *gin.Context) {\n\tm := models.{{SingularResource}}{}\n\tif err := c.ShouldBindUri(&m); err != nil {\n\t\tctrl.log.Error(\"invalid {{Resource}} deletion request\",\n\t\t\tzap.Error(err),\n\t\t)\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t\treturn\n\t}\n\t{{# ORM }}\n\t_, err := m.Delete(ctrl.db)\n\tif err != nil {\n\t\tctrl.log.Error(\"error deleting {{Resource}}\",\n\t\t\tzap.Error(err))\n\t\tc.AbortWithStatus(http.StatusInternalServerError)\n\t}\n\t{{/ ORM }}\n\tc.JSON(http.StatusNoContent, gin.H{})\n}\n"),
	}
	filee := &embedded.EmbeddedFile{
		Filename:    "gin/partials/index.tmpl",
		FileModTime: time.Unix(1594234732, 0),

		Content: string("// Index returns a list of {{Resource}} records\nfunc (ctrl *{{Resource}}Controller) {{Name}}(c *gin.Context) {\n\tvar results []models.{{SingularResource}}\n\t{{# ORM }}\n\tq := c.Request.URL.RawQuery\n\tqms := GetQueryModFromQuery(q)\n\tresults, err := models.{{PluralName}}(qms...).All(ctrl.db)\n\tif err != nil {\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t}\n\t{{/ ORM }}\n\tc.JSON(http.StatusOK, results)\n}\n"),
	}
	filef := &embedded.EmbeddedFile{
		Filename:    "gin/partials/show.tmpl",
		FileModTime: time.Unix(1594235682, 0),

		Content: string("// Show retrieves a new {{Resource}} record from the database\nfunc (ctrl *{{Resource}}Controller) {{Name}}(c *gin.Context) {\n\tvar result models.{{SingularResource}}\n\t{{# ORM }}\n\tid := c.GetInt(\"id\")\n\tresult, err := models.Find{{Resource}}(id)\n\tif err != nil {\n\t\tctrl.log.Error(\"error retrieving {{Resource}}\",\n\t\t\tzap.Error(err))\n\t\tc.AbortWithStatus(http.StatusInternalServerError)\n\t}\n\t{{/ ORM }}\n\tc.JSON(http.StatusOK, result)\n}\n"),
	}
	fileg := &embedded.EmbeddedFile{
		Filename:    "gin/partials/update.tmpl",
		FileModTime: time.Unix(1594234969, 0),

		Content: string("// Update updates a new {{Resource}} record in the database\nfunc (ctrl *{{Resource}}Controller) {{Name}}(c *gin.Context) {\n\tm := models.{{SingularResource}}{}\n\tif err := c.ShouldBindUri(&m); err != nil {\n\t\tctrl.log.Error(\"invalid {{Resource}} update request\",\n\t\t\tzap.Error(err),\n\t\t)\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t\treturn\n\t}\n\tif err := c.ShouldBindJSON(&m); err != nil {\n\t\tctrl.log.Error(\"invalid {{Resource}} update request\",\n\t\t\tzap.Error(err),\n\t\t)\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t\treturn\n\t}\n\t{{# ORM }}\n\terr := m.Update(ctrl.db, boil.Infer())\n\tif err != nil {\n\t\tctrl.log.Error(\"error updating {{Resource}}\",\n\t\t\tzap.Error(err))\n\t\tc.AbortWithStatus(http.StatusInternalServerError)\n\t}\n\t{{/ ORM }}\n\tc.JSON(http.StatusOK, gin.H{})\n}\n"),
	}
	fileh := &embedded.EmbeddedFile{
		Filename:    "gin/router.tpl",
		FileModTime: time.Unix(1594224538, 0),

		Content: string("package controllers\n\nimport (\n\t\"github.com/gin-gonic/gin\"\n\t\"github.com/jmoiron/sqlx\"\n\t\"go.uber.org/zap\"\n)\n\nfunc GetRouter(log *zap.SugaredLogger, db *sqlx.DB) *gin.Engine {\n\tr := gin.New()\n\n{{#Controllers}}\n\t{{Name}}Ctrl := {{Name}}Controller{db: db, log: log}\n{{#Operations}}\n\tr.{{Method}}(\"{{Path}}\", {{Name}}Ctrl.{{Handler}})\n{{/Operations}}\n{{/Controllers}}\n\treturn r\n}\n"),
	}
	filej := &embedded.EmbeddedFile{
		Filename:    "sqlboiler/query.go.tpl",
		FileModTime: time.Unix(1590423714, 0),

		Content: string("package controllers\n\nimport (\n\t\"fmt\"\n\t\"net/url\"\n\t\"strconv\"\n\n\t\"github.com/volatiletech/sqlboiler/queries/qm\"\n)\n\n// GetQueryModFromQuery derives db lookups from URI query parameters\nfunc GetQueryModFromQuery(query string) []qm.QueryMod {\n\tvar mods []qm.QueryMod\n\tm, _ := url.ParseQuery(query)\n\tfor k, v := range m {\n\t\tfor _, value := range v {\n\t\t\tif k == \"limit\" {\n\t\t\t\tlimit, err := strconv.Atoi(value)\n\t\t\t\tif err != nil {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\tmods = append(mods, qm.Limit(limit))\n\t\t\t} else if k == \"from\" {\n\t\t\t\tfrom, err := strconv.Atoi(value)\n\t\t\t\tif err != nil {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\t// TODO: support order by and ASC/DESC\n\t\t\t\tmods = append(mods, qm.Where(\"id >= ?\", from))\n\t\t\t} else {\n\t\t\t\tclause := fmt.Sprintf(\"%s=?\", k)\n\t\t\t\tmods = append(mods, qm.Where(clause, v))\n\t\t\t}\n\t\t}\n\t}\n\treturn mods\n}\n"),
	}
	filek := &embedded.EmbeddedFile{
		Filename:    "sqlboiler/seed_factory.tmpl",
		FileModTime: time.Unix(1590551496, 0),

		Content: string("package models\n\nimport (\n\t\"context\"\n\n\t\"github.com/bxcodec/faker\"\n\t\"github.com/jmoiron/sqlx\"\n\t\"github.com/volatiletech/sqlboiler/boil\"\n)\n\n// NewTest{{Name}} is a factory function to create fake/test data\nfunc NewTest{{Name}}() models.{{Name}} {\n  model := models.{{Name}}{}\n  faker.FakeData(&model)\n  return model\n}\n\n// Insert{{Name}} creates fake data for the {{Name}} model and inserts into the \n// database.\nfunc Insert{{Name}}(ctx context.Context, db *sqlx.DB, n int) error {\n  i := 0\n  for i < n {\n    m := NewTest{{Name}}()\n    if err := m.Insert(ctx, db, boil.Infer()); err != nil {\n      return err\n    }\n    i++\n  }\n\n  return nil\n}"),
	}
	filel := &embedded.EmbeddedFile{
		Filename:    "sqlboiler/seeder.tmpl",
		FileModTime: time.Unix(1590551496, 0),

		Content: string("package main\n\nimport (\n\t\"flag\"\n  \"seeds\"\n)\n\n// this seeder creates X number of records for each model\nfunc main() {\n\tconfigString := fmt.Sprintf(\"host=%s user=%s dbname=%s password=%s\", os.Getenv(\"POSTGRES_HOST\"), os.Getenv(\"POSTGRES_USER\"), os.Getenv(\"POSTGRES_DB\"), os.Getenv(\"POSTGRES_PASSWORD\"))\n\tdb, err := sqlx.Open(\"postgres\", configString)\n  if err != nil {\n    panic(err)\n  }\n\n  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n\tdefer cancel()\n\n  {{#each Models}}\n  if err := seed.Insert{{Name}}(ctx, db, 100); err != nil {\n    panic(err)\n  }\n  {{/each}}\n}\n"),
	}
	filen := &embedded.EmbeddedFile{
		Filename:    "tests/controller_test.tpl",
		FileModTime: time.Unix(1594236002, 0),

		Content: string("package controllers\n\nimport (\n\t\"net/http\"\n\t\"net/http/httptest\"\n\t\"testing\"\n\n\t\"{{ModuleName}}/models\"\n\n\t\"github.com/jmoiron/sqlx\"\n\tsqlmock \"github.com/DATA-DOG/go-sqlmock\"\n\t\"github.com/stretchr/testify/assert\"\n\t\"go.uber.org/zap\"\n)\n\n{{#each Actions}}\n{{{ whichActionTest Handler }}}\n{{/each}}\n"),
	}
	filep := &embedded.EmbeddedFile{
		Filename:    "tests/partials/create_test.tmpl",
		FileModTime: time.Unix(1594235941, 0),

		Content: string("func Test{{Resource}}Controller_{{Name}}(t *testing.T) {\n\tdb, mock, err := sqlmock.New()\n\tif err != nil {\n\t\tt.Fatalf(\"an error '%s' was not expected when opening a stub database connection: %s\", err, mock)\n\t}\n\tdefer db.Close()\n\tsqlxDB := sqlx.NewDb(db, \"sqlmock\")\n\n\ttests := []struct {\n\t\tname           string\n\t\tpath           string\n\t\twantStatusCode int\n\t}{\n\t\t{\n\t\t\tname:           \"Test creating with valid {{Resource}} as body\",\n\t\t\tpath:           \"{{Path}}\",\n\t\t\twantStatusCode: 201,\n\t\t},\n\t\t{\n\t\t\tname:           \"Test creating with empty request body\",\n\t\t\tpath:           \"{{Path}}\",\n\t\t\twantStatusCode: 400,\n\t\t},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\tlog := zap.NewExample().Sugar()\n\t\t\tdefer log.Sync()\n\n\t\t\trouter := GetRouter(log, sqlxDB)\n\n\t\t\tw := httptest.NewRecorder()\n\t\t\treq, _ := http.NewRequest(\"POST\", tt.path, nil)\n\t\t\trouter.ServeHTTP(w, req)\n\n\t\t\tassert.Equal(t, tt.wantStatusCode, w.Code)\n\t\t})\n\t}\n}\n"),
	}
	fileq := &embedded.EmbeddedFile{
		Filename:    "tests/partials/delete_test.tmpl",
		FileModTime: time.Unix(1594235966, 0),

		Content: string("func Test{{Resource}}Controller_{{Name}}(t *testing.T) {\n\tdb, mock, err := sqlmock.New()\n\tif err != nil {\n\t\tt.Fatalf(\"an error '%s' was not expected when opening a stub database connection: %s\", err, mock)\n\t}\n\tdefer db.Close()\n\tsqlxDB := sqlx.NewDb(db, \"sqlmock\")\n\n\ttests := []struct {\n\t\tname           string\n\t\tpath           string\n\t\twantStatusCode int\n\t}{\n\t\t{\n\t\t\tname:           \"Test deleting\",\n\t\t\tpath:           \"{{Path}}\",\n\t\t\twantStatusCode: 200,\n\t\t},\n\t\t{\n\t\t\tname:           \"Test deleting non-existent resource\",\n\t\t\tpath:           \"{{Path}}\",\n\t\t\twantStatusCode: 400,\n\t\t},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\tlog := zap.NewExample().Sugar()\n\t\t\tdefer log.Sync()\n\n\t\t\trouter := GetRouter(log, sqlxDB)\n\n\t\t\tw := httptest.NewRecorder()\n\t\t\treq, _ := http.NewRequest(\"DELETE\", tt.path, nil)\n\t\t\trouter.ServeHTTP(w, req)\n\n\t\t\tassert.Equal(t, tt.wantStatusCode, w.Code)\n\t\t})\n\t}\n}"),
	}
	filer := &embedded.EmbeddedFile{
		Filename:    "tests/partials/index_test.tmpl",
		FileModTime: time.Unix(1594235963, 0),

		Content: string("func Test{{Resource}}Controller_{{Name}}(t *testing.T) {\n\tdb, mock, err := sqlmock.New()\n\tif err != nil {\n\t\tt.Fatalf(\"an error '%s' was not expected when opening a stub database connection: %s\", err, mock)\n\t}\n\tdefer db.Close()\n\tsqlxDB := sqlx.NewDb(db, \"sqlmock\")\n\n\ttests := []struct {\n\t\tname           string\n\t\tpath           string\n\t\twant           []models.{{SingularResource}}\n\t\twantStatusCode int\n\t}{\n\t\t{\n\t\t\tname:           \"Test indexing without query parameters\",\n\t\t\tpath:           \"{{path}}\",\n\t\t\twant:           []models.{{SingularResource}}{},\n\t\t\twantStatusCode: 200,\n\t\t},\n\t\t{\n\t\t\tname:           \"Test indexing with parameters\",\n\t\t\tpath:           \"{{path}}?page=2\",\n\t\t\twant:           []models.{{SingularResource}}{},\n\t\t\twantStatusCode: 200,\n\t\t},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\tlog := zap.NewExample().Sugar()\n\t\t\tdefer log.Sync()\n\n\t\t\trouter := GetRouter(log, sqlxDB)\n\n\t\t\tw := httptest.NewRecorder()\n\t\t\treq, _ := http.NewRequest(\"GET\", tt.path, nil)\n\t\t\trouter.ServeHTTP(w, req)\n\n\t\t\tassert.Equal(t, tt.wantStatusCode, w.Code)\n\t\t\tassert.Equal(t, tt.want, w.Body.String())\n\t\t})\n\t}\n}\n"),
	}
	files := &embedded.EmbeddedFile{
		Filename:    "tests/partials/show_test.tmpl",
		FileModTime: time.Unix(1594236063, 0),

		Content: string("func Test{{Resource}}Controller_{{Name}}(t *testing.T) {\n\tdb, mock, err := sqlmock.New()\n\tif err != nil {\n\t\tt.Fatalf(\"an error '%s' was not expected when opening a stub database connection: %s\", err, mock)\n\t}\n\tdefer db.Close()\n\tsqlxDB := sqlx.NewDb(db, \"sqlmock\")\n\n  tests := []struct {\n    name           string\n    path           string\n    want           models.{{SingularResource}}\n    wantStatusCode int\n  }{\n    {\n      name:           \"Test getting existing {{Resource}}\",\n      path:           \"{{path}}\",\n      want:           models.{{SingularResource}}{},\n      wantStatusCode: 200,\n    },\n    {\n      name:           \"Test getting non-existent {{Resource}}\",\n      path:           \"{{path}}\",\n      want:           models.{{SingularResource}}{},\n      wantStatusCode: 200,\n    },\n  }\n  for _, tt := range tests {\n    t.Run(tt.name, func(t *testing.T) {\n\t\t\tlog := zap.NewExample().Sugar()\n\t\t\tdefer log.Sync()\n\n\t\t\trouter := GetRouter(log, sqlxDB)\n\n      w := httptest.NewRecorder()\n      req, _ := http.NewRequest(\"GET\", tt.path, nil)\n      router.ServeHTTP(w, req)\n\n      assert.Equal(t, tt.wantStatusCode, w.Code)\n      assert.Equal(t, tt.want, w.Body.String())\n    })\n  }\n}\n"),
	}
	filet := &embedded.EmbeddedFile{
		Filename:    "tests/partials/update_test.tmpl",
		FileModTime: time.Unix(1594235953, 0),

		Content: string("func Test{{Resource}}Controller_{{Name}}(t *testing.T) {\n\tdb, mock, err := sqlmock.New()\n\tif err != nil {\n\t\tt.Fatalf(\"an error '%s' was not expected when opening a stub database connection: %s\", err, mock)\n\t}\n\tdefer db.Close()\n\tsqlxDB := sqlx.NewDb(db, \"sqlmock\")\n\n\ttests := []struct {\n\t\tname           string\n\t\tpath           string\n\t\twantStatusCode int\n\t}{\n\t\t{\n\t\t\tname:           \"Test replacing with valid {{Resource}} as body\",\n\t\t\tpath:           \"{{path}}\",\n\t\t\twantStatusCode: 200,\n\t\t},\n\t\t{\n\t\t\tname:           \"Test replacing with empty request body\",\n\t\t\tpath:           \"{{path}}\",\n\t\t\twantStatusCode: 400,\n\t\t},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\tlog := zap.NewExample().Sugar()\n\t\t\tdefer log.Sync()\n\n\t\t\trouter := GetRouter(log, sqlxDB)\n\n\t\t\tw := httptest.NewRecorder()\n\t\t\treq, _ := http.NewRequest(\"PUT\", tt.path, nil)\n\t\t\trouter.ServeHTTP(w, req)\n\n\t\t\tassert.Equal(t, tt.wantStatusCode, w.Code)\n\t\t})\n\t}\n}\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1594235091, 0),
		ChildFiles: []*embedded.EmbeddedFile{},
	}
	dir2 := &embedded.EmbeddedDir{
		Filename:   "build",
		DirModTime: time.Unix(1587572809, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file3, // "build/circleciconfig.yml.tpl"
			file4, // "build/docker-compose.yml.tpl"
			file5, // "build/env.tpl"
			file6, // "build/sqlboiler.toml.tpl"
			file7, // "build/wait-for-server-start.sh.tpl"

		},
	}
	dir8 := &embedded.EmbeddedDir{
		Filename:   "gin",
		DirModTime: time.Unix(1592852870, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file9, // "gin/controller.tmpl"
			filea, // "gin/main.tpl"
			fileh, // "gin/router.tpl"

		},
	}
	dirb := &embedded.EmbeddedDir{
		Filename:   "gin/partials",
		DirModTime: time.Unix(1592852870, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			filec, // "gin/partials/create.tmpl"
			filed, // "gin/partials/delete.tmpl"
			filee, // "gin/partials/index.tmpl"
			filef, // "gin/partials/show.tmpl"
			fileg, // "gin/partials/update.tmpl"

		},
	}
	diri := &embedded.EmbeddedDir{
		Filename:   "sqlboiler",
		DirModTime: time.Unix(1592800172, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			filej, // "sqlboiler/query.go.tpl"
			filek, // "sqlboiler/seed_factory.tmpl"
			filel, // "sqlboiler/seeder.tmpl"

		},
	}
	dirm := &embedded.EmbeddedDir{
		Filename:   "tests",
		DirModTime: time.Unix(1591537292, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			filen, // "tests/controller_test.tpl"

		},
	}
	diro := &embedded.EmbeddedDir{
		Filename:   "tests/partials",
		DirModTime: time.Unix(1594232219, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			filep, // "tests/partials/create_test.tmpl"
			fileq, // "tests/partials/delete_test.tmpl"
			filer, // "tests/partials/index_test.tmpl"
			files, // "tests/partials/show_test.tmpl"
			filet, // "tests/partials/update_test.tmpl"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{
		dir2, // "build"
		dir8, // "gin"
		diri, // "sqlboiler"
		dirm, // "tests"

	}
	dir2.ChildDirs = []*embedded.EmbeddedDir{}
	dir8.ChildDirs = []*embedded.EmbeddedDir{
		dirb, // "gin/partials"

	}
	dirb.ChildDirs = []*embedded.EmbeddedDir{}
	diri.ChildDirs = []*embedded.EmbeddedDir{}
	dirm.ChildDirs = []*embedded.EmbeddedDir{
		diro, // "tests/partials"

	}
	diro.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`templates`, &embedded.EmbeddedBox{
		Name: `templates`,
		Time: time.Unix(1594235091, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"":               dir1,
			"build":          dir2,
			"gin":            dir8,
			"gin/partials":   dirb,
			"sqlboiler":      diri,
			"tests":          dirm,
			"tests/partials": diro,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"build/circleciconfig.yml.tpl":       file3,
			"build/docker-compose.yml.tpl":       file4,
			"build/env.tpl":                      file5,
			"build/sqlboiler.toml.tpl":           file6,
			"build/wait-for-server-start.sh.tpl": file7,
			"gin/controller.tmpl":                file9,
			"gin/main.tpl":                       filea,
			"gin/partials/create.tmpl":           filec,
			"gin/partials/delete.tmpl":           filed,
			"gin/partials/index.tmpl":            filee,
			"gin/partials/show.tmpl":             filef,
			"gin/partials/update.tmpl":           fileg,
			"gin/router.tpl":                     fileh,
			"sqlboiler/query.go.tpl":             filej,
			"sqlboiler/seed_factory.tmpl":        filek,
			"sqlboiler/seeder.tmpl":              filel,
			"tests/controller_test.tpl":          filen,
			"tests/partials/create_test.tmpl":    filep,
			"tests/partials/delete_test.tmpl":    fileq,
			"tests/partials/index_test.tmpl":     filer,
			"tests/partials/show_test.tmpl":      files,
			"tests/partials/update_test.tmpl":    filet,
		},
	})
}

func init() {

	// define files
	filev := &embedded.EmbeddedFile{
		Filename:    "Dockerfile",
		FileModTime: time.Unix(1587572809, 0),

		Content: string("FROM golang:1.13\n\nADD . /app\nWORKDIR /app\n\nCMD go run main.go"),
	}
	filew := &embedded.EmbeddedFile{
		Filename:    "Makefile",
		FileModTime: time.Unix(1590551496, 0),

		Content: string(".PHONY: models\n\n# Go parameters\nGOBUILD=go build\nGOCLEAN=go clean\nGOTEST=go test\nGOGET=go get\n\nall: test build\n\ndev-dependencies:\n\tgo get -u -t github.com/volatiletech/sqlboiler\n\tgo get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql\n\nbuild: \n\t$(GOBUILD) -tags=jsoniter .\n\ntest: \n\t$(GOTEST) -v ./...\n\nstart:\n\tmake build\n\tgo run main.go\n\n# usage: make migration N=tableName\nmigration:\n\tmigrate create -ext sql -dir ./migrations -seq $(N)\n\nmigratedb:\n\tmigrate up\n\ndropdb:\n\tmigrate drop\n\nmodels:\n\tsqlboiler psql --no-tests --no-hooks --no-context\n"),
	}
	filex := &embedded.EmbeddedFile{
		Filename:    "test.Dockerfile",
		FileModTime: time.Unix(1587572809, 0),

		Content: string("FROM golang:1.13\n\nRUN go get -u github.com/smartystreets/goconvey\n\nADD . /app\nWORKDIR /app\n\nRUN go install -v\n\nCMD goconvey -host 0.0.0.0 -port=9999\n\nEXPOSE 9999"),
	}

	// define dirs
	diru := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1590551496, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			filev, // "Dockerfile"
			filew, // "Makefile"
			filex, // "test.Dockerfile"

		},
	}

	// link ChildDirs
	diru.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`static`, &embedded.EmbeddedBox{
		Name: `static`,
		Time: time.Unix(1590551496, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": diru,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"Dockerfile":      filev,
			"Makefile":        filew,
			"test.Dockerfile": filex,
		},
	})
}
