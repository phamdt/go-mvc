package gomvc

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/spf13/cobra"
)

var oa = &cobra.Command{
	Use:   "oa",
	Short: "Generate controllers from an OpenAPI yml file",
	Run: func(cmd *cobra.Command, args []string) {
		configDir, err := cmd.LocalFlags().GetString("config")
		if err != nil {
			log.Println(err.Error())
			return
		}
		// TODO: read spec location from config
		specPath, err := cmd.LocalFlags().GetString("spec")
		if err != nil {
			log.Println(err.Error())
			return
		}

		// read intended destination for generation output
		destDir, err := cmd.LocalFlags().GetString("dest")
		if err != nil {
			log.Println(err.Error())
			return
		}

		templateDir, err := cmd.LocalFlags().GetString("templates")
		if err != nil {
			log.Println(err.Error())
			return
		}
		oa3 := LoadWithKin(specPath)
		GenerateFromOA(oa3, destDir, templateDir, configDir)
	},
}

func GenerateFromOA(oa3 *openapi3.Swagger, dest, templateDir, configDir string) {
	config := NewGoMVCConfig(configDir)

	createDirIfNotExists(dest)
	ctrlDest := filepath.Join(dest, "controllers")
	createDirIfNotExists(ctrlDest)

	CreateRouter(RouteData{}, "gin/router.tpl", ctrlDest)
	g := NewOAGenerator(oa3)
	for path, pathItem := range oa3.Paths {
		path = strings.Trim(path, " ")
		log.Printf("examining path, %s\n", path)
		if config.IsBlacklisted(path) {
			continue
		}
		if err := g.CreateControllerFiles(path, pathItem, dest, templateDir); err != nil {
			log.Fatalf("%s: %s", path, err.Error())
		}
	}
}

// OAGenerator wraps functionality for reading and manipulating a single OpenAPI spec
type OAGenerator struct {
	oa3 *openapi3.Swagger
}

// NewOAGenerator is a constructor for OAGenerator
func NewOAGenerator(oa3 *openapi3.Swagger) OAGenerator {
	return OAGenerator{
		oa3: oa3,
	}
}

// CreateControllerFiles creates a controller for operations found in
// an OpenAPI file
func (oag *OAGenerator) CreateControllerFiles(path string, pathItem *openapi3.PathItem, dest string, templateDir string) error {
	name := strcase.ToSnake(pathItem.Summary)
	name = strings.ToLower(name)

	// skip creating file if we can't find a good name from the doc
	if name == "" {
		log.Printf("No summary provided in API, defaulting to deriving name from path %s since we can't identify a name for the resource", path)
		name = getControllerNameFromPath(path)
	}
	log.Printf("Preparing to generate controller files for %s %s", name, path)
	data := ControllerData{
		Name:       name,
		PluralName: inflection.Plural(name),
		Path:       path,
		Actions:    []Action{},
	}

	// collect controller methods based on specified HTTP verbs/operations
	for method, op := range pathItem.Operations() {
		var handler = getDefaultHandlerName(method, path)
		var operationName string
		if op.OperationID == "" {
			log.Printf("Missing operation ID. Generating default name for handler/operation function in controller %s.\n", name)
			operationName = handler
		} else {
			operationName = strings.Title(op.OperationID)
		}

		action := Action{
			Method: method, Path: path, Handler: handler, Name: operationName,
			Resource: name,
		}
		data.Actions = append(data.Actions, action)
	}

	if err := createControllerFromDefault(data, dest); err != nil {
		return err
	}
	log.Printf("created controller actions for %s\n", path)
	return nil
}

var methodLookup = map[string]string{
	"GET":    "Show",
	"POST":   "Create",
	"PUT":    "Update",
	"DELETE": "Delete",
}

func getDefaultHandlerName(method, path string) string {
	var handler string
	// index is a special case because currently, all of the HTTP verbs have a 1
	// to 1 relationship with an action except for Index and Show which both use
	// GET.
	if isIndex(method, path) {
		handler = "Index"
	} else {
		handler = methodLookup[method]
	}
	return handler
}

// OA is the cli command that creates a router and controller functions from an
// OpenAPI file
func OA() *cobra.Command {
	return oa
}
