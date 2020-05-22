package gomvc

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "Dockerfile",
		FileModTime: time.Unix(1582525394, 0),

		Content: string("FROM golang:1.13\n\nADD . /app\nWORKDIR /app\n\nCMD go run main.go"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "Makefile",
		FileModTime: time.Unix(1582431832, 0),

		Content: string("# Go parameters\nGOBUILD=go build\nGOCLEAN=go clean\nGOTEST=go test\nGOGET=go get\n\nall: test build\n\ndev-dependencies:\n\tgo get -u -t github.com/volatiletech/sqlboiler\n\tgo get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql\n\nbuild: \n\t$(GOBUILD) -tags=jsoniter .\ntest: \n\t$(GOTEST) -v ./...\nstart:\n\tgo build .\n\tgo run main.go"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "test.Dockerfile",
		FileModTime: time.Unix(1582525396, 0),

		Content: string("FROM golang:1.13\n\nRUN go get -u github.com/smartystreets/goconvey\n\nADD . /app\nWORKDIR /app\n\nRUN go install -v\n\nCMD goconvey -host 0.0.0.0 -port=9999\n\nEXPOSE 9999"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1583943616, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "Dockerfile"
			file3, // "Makefile"
			file4, // "test.Dockerfile"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`static`, &embedded.EmbeddedBox{
		Name: `static`,
		Time: time.Unix(1583943616, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"Dockerfile":      file2,
			"Makefile":        file3,
			"test.Dockerfile": file4,
		},
	})
}

func init() {

	// define files
	file6 := &embedded.EmbeddedFile{
		Filename:    "README.md",
		FileModTime: time.Unix(1584700131, 0),

		Content: string(""),
	}
	file8 := &embedded.EmbeddedFile{
		Filename:    "build/circleciconfig.yml.tpl",
		FileModTime: time.Unix(1582431518, 0),

		Content: string("version: 2\njobs:\n  build_and_test:\n    docker:\n      - image: circleci/golang:1.13\n    working_directory: /go/src/{{gitRepoPath}}\n    steps:\n      - checkout\n      - setup_remote_docker:\n          docker_layer_caching: true\n      - add_ssh_keys\n{{#envFileName}}\n      - run:\n          name: Add environment variables to a file\n          command: cp {{#envFileSampleName}} {{envFileName}}\n{{/envFileName}}\n      - run:\n          name: Start Containers\n          command: docker-compose -f docker-compose.yml up -d\n      - run:\n          name: Wait for Server\n          command: |\n            chmod +x .circleci/wait-for-server-start.sh\n            .circleci/wait-for-server-start.sh\n      - run:\n          name: Wait extra 10s to ensure database is seeded\n          command: sleep 10\n      - run:\n          name: Run tests\n          command: docker exec -it {{containerName}} go test ./...\n\nworkflows:\n  version: 2\n  build:\n    jobs:\n      - build_and_test"),
	}
	file9 := &embedded.EmbeddedFile{
		Filename:    "build/docker-compose.yml.tpl",
		FileModTime: time.Unix(1583864636, 0),

		Content: string("version: \"3\"\nservices:\n  {{Name}}_postgres:\n    container_name: {{Name}}_db\n    hostname: {{Name}}_db\n    image: \"postgres:11\"\n    env_file: .env\n    ports:\n      - \"5432:5432\"\n# UNCOMMENT ONCE YOU HAVE MIGRATIONS\n#  {{Name}}_migrations:\n#    container_name: migrations\n#    image: migrate/migrate:v4.6.2\n#    command: [\"-path\", \"/migrations/\", \"-database\", \"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable\", \"up\"]\n#    depends_on:\n#      - postgres\n#    env_file: .env\n#    restart: on-failure\n#    links: \n#      - postgres\n#    volumes:\n#      - ./migrations:/migrations \n#\n  {{Name}}:\n    container_name: {{Name}}\n    build:\n      context: .\n      dockerfile: Dockerfile\n    env_file: .env\n    volumes:\n      - ./:/go/src/{{Name}}\n    ports:\n      - \"8080:8080\"\n    links:\n      - {{Name}}_postgres\n\n  {{Name}}_test:\n    container_name: {{Name}}_test\n    build:\n      context: .\n      dockerfile: test.Dockerfile\n    env_file: .env\n    volumes:\n      - ./:/go/src/{{Name}}\n    ports:\n      - \"9999:9999\"\n    links:\n      - {{Name}}_postgres\n\n"),
	}
	filea := &embedded.EmbeddedFile{
		Filename:    "build/env.tpl",
		FileModTime: time.Unix(1583690742, 0),

		Content: string("# Postgres Database\n# Env vars originate from the postgres image on dockerhub\nPOSTGRES_HOST={{Name}}\nPOSTGRES_USER={{Name}}\nPOSTGRES_DB={{Name}}\nPOSTGRES_PASSWORD={{Name}}\n\nAPP_NAME={{Name}}\nNR_LICENSE_KEY="),
	}
	fileb := &embedded.EmbeddedFile{
		Filename:    "build/sqlboiler.toml.tpl",
		FileModTime: time.Unix(1583913888, 0),

		Content: string("[psql]\n  dbname = \"{{dbName}}\"\n  host   = \"0.0.0.0\"\n  port   = {{dbPort}}\n  user   = \"{{dbUser}}\"\n  pass   = \"{{dbPassword}}\"\n  blacklist = [\n    {{#blacklist}}{{blacklist}}{{/blacklist}}\n  ]\n  sslmode = \"disable\"\n{{#templates}}\n  templates = [\n    {{templates}}\n  ]\n{{/templates}}\n"),
	}
	filec := &embedded.EmbeddedFile{
		Filename:    "build/wait-for-server-start.sh.tpl",
		FileModTime: time.Unix(1583864824, 0),

		Content: string("#!/bin/bash\n\necho \"Waiting for servers to start...\"\nattempts=1\nwhile true; do\n  docker exec -i {{Name}} curl -f http://localhost:8080/health > /dev/null 2> /dev/null\n  if [ $? = 0 ]; then\n    echo \"Service started\"\n    break\n  fi\n  ((attempts++))\n  if [[ $attempts == 5 ]]; then\n    echo \"5 attempts to check health failed\"\n    break\n  fi\n  sleep 10\n  echo $attempts\ndone"),
	}
	filee := &embedded.EmbeddedFile{
		Filename:    "gin/controller.gotmpl",
		FileModTime: time.Unix(1580885544, 0),

		Content: string("// Code generated by ***REMOVED***swagger; DO NOT EDIT.\n{{ if .Copyright -}}// {{ comment .Copyright -}}{{ end }}\n\npackage controller\n\nimport (\n  \"bytes\"\n  \"encoding/json\"\n  \"github.com/gin-gonic/gin\"\n  \"io/ioutil\"\n  \"log\"\n  \"net/http\"\n\n  {{ range $key, $value := .Imports }}{{ $key }} {{ printf \"%q\" $value }}\n  {{ end }}\n)\n\nfunc {{ pascalize .Name }}(c *gin.Context) {\n  {{ if .SuccessResponse.IsSuccess }}var data {{ pascalize .SuccessResponse.Name }}\n  {{ else }}var data string\n  {{ end }}\n\n  errText := \"\"\n  resp, err := Request(\"{{ .Method }}\", path{{ pascalize .Name }}(c), nil)\n  defer resp.Body.Close()\n  if err != nil {\n    errText = err.Error()\n  } else {\n    result, err := ioutil.ReadAll(resp.Body)\n    err = json.Unmarshal(result, &data)\n    if err != nil {\n      errText = err.Error()\n    }\n  }\n  c.HTML(http.StatusOK, \"{{ snakize .Name }}.html\", gin.H{\n    \"title\": \"{{ .Summary }}\",\n    \"data\": data,\n    \"error\": errText,\n    \"type\": \"{{ .SuccessResponse.Name }}\",\n  })\n}\n\nfunc path{{ pascalize .Name }}(c *gin.Context) string {\n  path := \"{{ .BasePath }}{{ .Path }}\"\n  {{if .PathParams }}\n    replaces := []string{\n    {{ range $key, $value := .PathParams }}{{ $value.Path }},{{ end }}\n    }\n    for i := range replaces {\n    path = strings.Replace(path, \"{\" + replaces[i] + \"}\", c.Param(replaces[i]), 1)\n    }\n  {{ end }}\n  return path\n}"),
	}
	filef := &embedded.EmbeddedFile{
		Filename:    "gin/controller.tmpl",
		FileModTime: time.Unix(1582614378, 0),

		Content: string("package controllers\n\nimport (\n\t\"net/http\"\n\n\t\"github.com/gin-gonic/gin\"\n\t\"github.com/jmoiron/sqlx\"\n\t\"github.com/volatiletech/sqlboiler/boil\"\n\t\"go.uber.org/zap\"\n)\n\n// {{Name}}Controller exposes the methods for interacting with the\n// RESTful {{Name}} resource\ntype {{Name}}Controller struct {\n\tdb  *sqlx.DB\n\tlog *zap.Logger\n}\n\n{{#each Actions}}\n{{{ whichAction Name }}}\n{{/each}}"),
	}
	fileg := &embedded.EmbeddedFile{
		Filename:    "gin/main.tpl",
		FileModTime: time.Unix(1583690768, 0),

		Content: string("package main\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"log\"\n\t\"net/http\"\n\t\"os\"\n\t\"os/signal\"\n\t\"syscall\"\n\t\"time\"\n\t\"{{Name}}/controllers\"\n\n\t\"github.com/gin-gonic/gin\"\n\t\"github.com/jmoiron/sqlx\"\n\t_ \"github.com/lib/pq\" // blank import necessary to use driver\n\tnewrelic \"github.com/newrelic/go-agent\"\n\t\"github.com/newrelic/go-agent/_integrations/nrgin/v1\"\n\t\"go.uber.org/zap\"\n)\n\nfunc main() {\n\t// construct dependencies\n\tlog := zap.NewExample().Sugar()\n\tdefer log.Sync()\n\n\t// setup database\n\tdb, err := newDb()\n\tif err != nil {\n\t\tlog.Fatalf(\"can't initalize database connection: %v\", zap.Error(err))\n\t\treturn\n\t}\n\n\t// setup router and middleware\n\trouter := controllers.GetRouter(log, db)\n\t// Recovery middleware recovers from any panics and writes a 500 if there was one.\n\trouter.Use(gin.Recovery())\n\n\t// setup monitoring only if the license key is set\n\tnrKey := os.Getenv(\"NR_LICENSE_KEY\")\n\tif nrKey != \"\" {\n\t\tnrMiddleware, err := newRelic(nrKey)\n\t\tif err != nil {\n\t\t\tlog.Fatal(\"Unexpected error setting up new relic\", zap.Error(err))\n\t\t\tpanic(err)\n\t\t}\n\t\trouter.Use(nrMiddleware)\n\t}\n\n\tsrv := &http.Server{\n\t\tAddr:    \":8080\",\n\t\tHandler: router,\n\t}\n\n\tgo func() {\n\t\t// service connections\n\t\tif err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {\n\t\t\tlog.Fatalf(\"listen: %s\\n\", zap.Error(err))\n\t\t}\n\t}()\n\n\t// Wait for interrupt signal to gracefully shutdown the server with\n\t// a timeout of 5 seconds.\n\tquit := make(chan os.Signal)\n\t// kill (no param) default send syscall.SIGTERM\n\t// kill -2 is syscall.SIGINT\n\t// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it\n\tsignal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)\n\t<-quit\n\tlog.Info(\"Shutdown Server ...\")\n\n\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n\tdefer cancel()\n\tif err := srv.Shutdown(ctx); err != nil {\n\t\tlog.Fatal(\"Server Shutdown:\", zap.Error(err))\n\t}\n\t// catching ctx.Done(). timeout of 5 seconds.\n\tselect {\n\tcase <-ctx.Done():\n\t\tlog.Info(\"timeout of 5 seconds.\")\n\t}\n\tlog.Info(\"Server exiting\")\n}\n\nfunc newRelic(nrKey string) (gin.HandlerFunc, error) {\n\tcfg := newrelic.NewConfig(os.Getenv(\"APP_NAME\"), nrKey)\n\t// Creates a New Relic Application\n\tapm, err := newrelic.NewApplication(cfg)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn nrgin.Middleware(apm), nil\n}\n\nfunc newDb() (*sqlx.DB, error) {\n\tconfigString := fmt.Sprintf(\"host=%s user=%s dbname=%s password=%s\", os.Getenv(\"POSTGRES_HOST\"), os.Getenv(\"POSTGRES_USER\"), os.Getenv(\"POSTGRES_DB\"), os.Getenv(\"POSTGRES_PASSWORD\"))\n\treturn sqlx.Open(\"postgres\", configString)\n}\n"),
	}
	filei := &embedded.EmbeddedFile{
		Filename:    "gin/partials/create.tmpl",
		FileModTime: time.Unix(1582614262, 0),

		Content: string("// Create saves a new {{Name}} record into the database\nfunc (ctrl *{{Name}}Controller) Create(c *gin.Context) {\n\tm := models.{{Name}}{}\n\tif err := c.ShouldBindJSON(m); err != nil {\n\t\tctrl.log.Error(\"invalid {{Name}} creation request\",\n\t\t\tzap.Error(err),\n\t\t)\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t\treturn\n\t}\n\terr := m.Insert(ctrl.db, boil.Infer())\n\tif err != nil {\n\t\tctrl.log.Error(\"error creating {{Name}}\",\n\t\t\tzap.Error(err))\n\t\tc.AbortWithStatus(http.StatusInternalServerError)\n\t}\n\tc.JSON(http.StatusCreated, gin.H{})\n}\n"),
	}
	filej := &embedded.EmbeddedFile{
		Filename:    "gin/partials/delete.tmpl",
		FileModTime: time.Unix(1582613836, 0),

		Content: string("// Delete deletes a new {{Name}} record into the database\nfunc (ctrl *{{Name}}Controller) Delete(c *gin.Context) {\n\tm := models.{{Name}}{}\n\tif err := c.ShouldBindUri(&m); err != nil {\n\t\tctrl.log.Error(\"invalid {{Name}} deletion request\",\n\t\t\tzap.Error(err),\n\t\t)\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t\treturn\n\t}\n\terr := m.Delete(ctrl.db)\n\tif err != nil {\n\t\tctrl.log.Error(\"error deleting {{Name}}\",\n\t\t\tzap.Error(err))\n\t\tc.AbortWithStatus(http.StatusInternalServerError)\n\t}\n\tc.JSON(http.StatusOK, gin.H{})\n}\n"),
	}
	filek := &embedded.EmbeddedFile{
		Filename:    "gin/partials/index.tmpl",
		FileModTime: time.Unix(1582613944, 0),

		Content: string("// Index returns a list of {{Name}} records\nfunc (ctrl *{{Name}}Controller) Index(c *gin.Context) {\n\tq := c.Request.URL.RawQuery\n\tqms := GetQueryModFromQuery(q)\n\tresults, err := models.{{Name}}(qms...).All(ctrl.db)\n\tif err != nil {\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t}\n\tc.JSON(http.StatusOK, results)\n}\n"),
	}
	filel := &embedded.EmbeddedFile{
		Filename:    "gin/partials/show.tmpl",
		FileModTime: time.Unix(1582588204, 0),

		Content: string("// Show retrieves a new {{Name}} record from the database\nfunc (ctrl *{{Name}}Controller) Show(c *gin.Context) {\n\tm := models.{{Name}}{}\n\tif err := c.ShouldBindUri(&m); err != nil {\n\t\tctrl.log.Error(\"invalid {{Name}} retrieval request\",\n\t\t\tzap.Error(err),\n\t\t)\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t\treturn\n\t}\n\tresult, err := models.Find{{Name}}(id)\n\tif err != nil {\n\t\tctrl.log.Error(\"error retrieving {{Name}}\",\n\t\t\tzap.Error(err))\n\t\tc.AbortWithStatus(http.StatusInternalServerError)\n\t}\n\tc.JSON(http.StatusOK, result)\n}\n"),
	}
	filem := &embedded.EmbeddedFile{
		Filename:    "gin/partials/update.tmpl",
		FileModTime: time.Unix(1582588234, 0),

		Content: string("// Update updates a new {{Name}} record in the database\nfunc (ctrl *{{Name}}Controller) Update(c *gin.Context) {\n\tm := models.{{Name}}{}\n\tif err := c.ShouldBindUri(&m); err != nil {\n\t\tctrl.log.Error(\"invalid {{Name}} update request\",\n\t\t\tzap.Error(err),\n\t\t)\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t\treturn\n\t}\n\tif err := c.ShouldBindJSON(&m); err != nil {\n\t\tctrl.log.Error(\"invalid {{Name}} update request\",\n\t\t\tzap.Error(err),\n\t\t)\n\t\tc.AbortWithError(http.StatusBadRequest, err)\n\t\treturn\n\t}\n\terr := m.Update(ctrl.db, boil.Infer())\n\tif err != nil {\n\t\tctrl.log.Error(\"error updating {{Name}}\",\n\t\t\tzap.Error(err))\n\t\tc.AbortWithStatus(http.StatusInternalServerError)\n\t}\n\tc.JSON(http.StatusOK, gin.H{})\n}\n"),
	}
	filen := &embedded.EmbeddedFile{
		Filename:    "gin/router.tpl",
		FileModTime: time.Unix(1582764064, 0),

		Content: string("package controllers\n\nimport (\n\t\"github.com/gin-gonic/gin\"\n\t\"github.com/jmoiron/sqlx\"\n\t\"go.uber.org/zap\"\n)\n\nfunc GetRouter(log *zap.SugaredLogger, db *sqlx.DB) *gin.Engine {\n\tr := gin.New()\n\n{{#Controllers}}\n\t{{Name}}Ctrl := {{Name}}Controller{}\n{{#Operations}}\n\tr.{{Method}}(\"{{Path}}\", {{Name}}Ctrl.{{Handler}})\n{{/Operations}}\n{{/Controllers}}\n\treturn r\n}\n"),
	}
	fileo := &embedded.EmbeddedFile{
		Filename:    "query.go.tpl",
		FileModTime: time.Unix(1583961962, 0),

		Content: string("package controllers\n\nimport (\n\t\"fmt\"\n\t\"net/url\"\n\t\"strconv\"\n\n\t\"github.com/volatiletech/sqlboiler/queries/qm\"\n)\n\n// GetQueryModFromQuery derives db lookups from URI query parameters\nfunc GetQueryModFromQuery(query string) []qm.QueryMod {\n\tvar mods []qm.QueryMod\n\tm, _ := url.ParseQuery(query)\n\tfor k, v := range m {\n\t\tfor _, value := range v {\n\t\t\tif k == \"limit\" {\n\t\t\t\tlimit, err := strconv.Atoi(value)\n\t\t\t\tif err != nil {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\tmods = append(mods, qm.Limit(limit))\n\t\t\t} else if k == \"from\" {\n\t\t\t\tfrom, err := strconv.Atoi(value)\n\t\t\t\tif err != nil {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\t// TODO: support order by and ASC/DESC\n\t\t\t\tmods = append(mods, qm.Where(\"id >= ?\", from))\n\t\t\t} else {\n\t\t\t\tclause := fmt.Sprintf(\"%s=?\", k)\n\t\t\t\tmods = append(mods, qm.Where(clause, v))\n\t\t\t}\n\t\t}\n\t}\n\treturn mods\n}\n"),
	}
	fileq := &embedded.EmbeddedFile{
		Filename:    "tests/controller_test.tpl",
		FileModTime: time.Unix(1582615456, 0),

		Content: string("package controllers\n\nimport (\n\t\"net/http\"\n\t\"net/http/httptest\"\n\t\"testing\"\n\n\t\"github.com/stretchr/testify/assert\"\n)\n\n{{#each Actions}}\n{{{ whichActionTest Name }}}\n{{/each}}\n"),
	}
	files := &embedded.EmbeddedFile{
		Filename:    "tests/partials/create_test.tmpl",
		FileModTime: time.Unix(1582615488, 0),

		Content: string("func Test{{Name}}Controller_Create(t *testing.T) {\n\ttests := []struct {\n\t\tname           string\n\t\tpath           string\n\t\twantStatusCode int\n\t}{\n\t\t{\n\t\t\tname:           \"Test creating with valid {{Name}} as body\",\n\t\t\tpath:           \"{{Path}}\",\n\t\t\twantStatusCode: 201,\n\t\t},\n\t\t{\n\t\t\tname:           \"Test creating with empty request body\",\n\t\t\tpath:           \"{{Path}}\",\n\t\t\twantStatusCode: 400,\n\t\t},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\trouter := GetRouter()\n\n\t\t\tw := httptest.NewRecorder()\n\t\t\treq, _ := http.NewRequest(\"POST\", tt.path, nil)\n\t\t\trouter.ServeHTTP(w, req)\n\n\t\t\tassert.Equal(t, tt.wantStatusCode, w.Code)\n\t\t})\n\t}\n}\n"),
	}
	filet := &embedded.EmbeddedFile{
		Filename:    "tests/partials/delete_test.tmpl",
		FileModTime: time.Unix(1582701988, 0),

		Content: string("func Test{{Name}}Controller_Delete(t *testing.T) {\n\ttests := []struct {\n\t\tname           string\n\t\tpath           string\n\t\twantStatusCode int\n\t}{\n\t\t{\n\t\t\tname:           \"Test deleting\",\n\t\t\tpath:           \"{{Path}}\",\n\t\t\twantStatusCode: 200,\n\t\t},\n\t\t{\n\t\t\tname:           \"Test deleting non-existent resource\",\n\t\t\tpath:           \"{{Path}}\",\n\t\t\twantStatusCode: 400,\n\t\t},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\trouter := GetRouter()\n\n\t\t\tw := httptest.NewRecorder()\n\t\t\treq, _ := http.NewRequest(\"DELETE\", tt.path, nil)\n\t\t\trouter.ServeHTTP(w, req)\n\n\t\t\tassert.Equal(t, tt.wantStatusCode, w.Code)\n\t\t})\n\t}\n}"),
	}
	fileu := &embedded.EmbeddedFile{
		Filename:    "tests/partials/index_test.tmpl",
		FileModTime: time.Unix(1582615488, 0),

		Content: string("func Test{{Name}}Controller_Index(t *testing.T) {\n\ttests := []struct {\n\t\tname           string\n\t\tpath           string\n\t\twant           []{{Name}}\n\t\twantStatusCode int\n\t}{\n\t\t{\n\t\t\tname:           \"Test indexing without query parameters\",\n\t\t\tpath:           \"{{path}}\",\n\t\t\twant:           []{{Name}}{},\n\t\t\twantStatusCode: 200,\n\t\t},\n\t\t{\n\t\t\tname:           \"Test indexing with parameters\",\n\t\t\tpath:           \"{{path}}?page=2\",\n\t\t\twant:           []{{Name}}{},\n\t\t\twantStatusCode: 200,\n\t\t},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\trouter := GetRouter()\n\n\t\t\tw := httptest.NewRecorder()\n\t\t\treq, _ := http.NewRequest(\"GET\", tt.path, nil)\n\t\t\trouter.ServeHTTP(w, req)\n\n\t\t\tassert.Equal(t, tt.wantStatusCode, w.Code)\n\t\t\tassert.Equal(t, tt.want, w.Body.String())\n\t\t})\n\t}\n}\n"),
	}
	filev := &embedded.EmbeddedFile{
		Filename:    "tests/partials/show_test.tmpl",
		FileModTime: time.Unix(1590000447, 0),

		Content: string("func Test{{Name}}Controller_Show(t *testing.T) {\n  tests := []struct {\n    name           string\n    path           string\n    want           []{{Name}}\n    wantStatusCode int\n  }{\n    {\n      name:           \"Test getting existing {{Name}}\",\n      path:           \"{{path}}\",\n      want:           []{{Name}}{},\n      wantStatusCode: 200,\n    },\n    {\n      name:           \"Test getting non-existent {{Name}}\",\n      path:           \"{{path}}\",\n      want:           []{{Name}}{},\n      wantStatusCode: 200,\n    },\n  }\n  for _, tt := range tests {\n    t.Run(tt.name, func(t *testing.T) {\n      router := GetRouter()\n\n      w := httptest.NewRecorder()\n      req, _ := http.NewRequest(\"GET\", tt.path, nil)\n      router.ServeHTTP(w, req)\n\n      assert.Equal(t, tt.wantStatusCode, w.Code)\n      assert.Equal(t, tt.want, w.Body.String())\n    })\n  }\n}\n"),
	}
	filew := &embedded.EmbeddedFile{
		Filename:    "tests/partials/update_test.tmpl",
		FileModTime: time.Unix(1582615488, 0),

		Content: string("func Test{{name}}Controller_Replace(t *testing.T) {\n\ttests := []struct {\n\t\tname           string\n\t\tpath           string\n\t\twantStatusCode int\n\t}{\n\t\t{\n\t\t\tname:           \"Test replacing with valid {{name}} as body\",\n\t\t\tpath:           \"{{path}}\",\n\t\t\twantStatusCode: 200,\n\t\t},\n\t\t{\n\t\t\tname:           \"Test replacing with empty request body\",\n\t\t\tpath:           \"{{path}}\",\n\t\t\twantStatusCode: 400,\n\t\t},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\trouter := GetRouter()\n\n\t\t\tw := httptest.NewRecorder()\n\t\t\treq, _ := http.NewRequest(\"PUT\", tt.path, nil)\n\t\t\trouter.ServeHTTP(w, req)\n\n\t\t\tassert.Equal(t, tt.wantStatusCode, w.Code)\n\t\t})\n\t}\n}\n"),
	}

	// define dirs
	dir5 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1589993614, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file6, // "README.md"
			fileo, // "query.go.tpl"

		},
	}
	dir7 := &embedded.EmbeddedDir{
		Filename:   "build",
		DirModTime: time.Unix(1583943618, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file8, // "build/circleciconfig.yml.tpl"
			file9, // "build/docker-compose.yml.tpl"
			filea, // "build/env.tpl"
			fileb, // "build/sqlboiler.toml.tpl"
			filec, // "build/wait-for-server-start.sh.tpl"

		},
	}
	dird := &embedded.EmbeddedDir{
		Filename:   "gin",
		DirModTime: time.Unix(1589641349, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			filee, // "gin/controller.gotmpl"
			filef, // "gin/controller.tmpl"
			fileg, // "gin/main.tpl"
			filen, // "gin/router.tpl"

		},
	}
	dirh := &embedded.EmbeddedDir{
		Filename:   "gin/partials",
		DirModTime: time.Unix(1583943618, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			filei, // "gin/partials/create.tmpl"
			filej, // "gin/partials/delete.tmpl"
			filek, // "gin/partials/index.tmpl"
			filel, // "gin/partials/show.tmpl"
			filem, // "gin/partials/update.tmpl"

		},
	}
	dirp := &embedded.EmbeddedDir{
		Filename:   "tests",
		DirModTime: time.Unix(1589993627, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			fileq, // "tests/controller_test.tpl"

		},
	}
	dirr := &embedded.EmbeddedDir{
		Filename:   "tests/partials",
		DirModTime: time.Unix(1583943618, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			files, // "tests/partials/create_test.tmpl"
			filet, // "tests/partials/delete_test.tmpl"
			fileu, // "tests/partials/index_test.tmpl"
			filev, // "tests/partials/show_test.tmpl"
			filew, // "tests/partials/update_test.tmpl"

		},
	}

	// link ChildDirs
	dir5.ChildDirs = []*embedded.EmbeddedDir{
		dir7, // "build"
		dird, // "gin"
		dirp, // "tests"

	}
	dir7.ChildDirs = []*embedded.EmbeddedDir{}
	dird.ChildDirs = []*embedded.EmbeddedDir{
		dirh, // "gin/partials"

	}
	dirh.ChildDirs = []*embedded.EmbeddedDir{}
	dirp.ChildDirs = []*embedded.EmbeddedDir{
		dirr, // "tests/partials"

	}
	dirr.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`templates`, &embedded.EmbeddedBox{
		Name: `templates`,
		Time: time.Unix(1589993614, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"":               dir5,
			"build":          dir7,
			"gin":            dird,
			"gin/partials":   dirh,
			"tests":          dirp,
			"tests/partials": dirr,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"README.md":                          file6,
			"build/circleciconfig.yml.tpl":       file8,
			"build/docker-compose.yml.tpl":       file9,
			"build/env.tpl":                      filea,
			"build/sqlboiler.toml.tpl":           fileb,
			"build/wait-for-server-start.sh.tpl": filec,
			"gin/controller.gotmpl":              filee,
			"gin/controller.tmpl":                filef,
			"gin/main.tpl":                       fileg,
			"gin/partials/create.tmpl":           filei,
			"gin/partials/delete.tmpl":           filej,
			"gin/partials/index.tmpl":            filek,
			"gin/partials/show.tmpl":             filel,
			"gin/partials/update.tmpl":           filem,
			"gin/router.tpl":                     filen,
			"query.go.tpl":                       fileo,
			"tests/controller_test.tpl":          fileq,
			"tests/partials/create_test.tmpl":    files,
			"tests/partials/delete_test.tmpl":    filet,
			"tests/partials/index_test.tmpl":     fileu,
			"tests/partials/show_test.tmpl":      filev,
			"tests/partials/update_test.tmpl":    filew,
		},
	})
}
