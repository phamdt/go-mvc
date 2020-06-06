package gomvc

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type ControllerData struct {
	Name       string
	PluralName string
	Path       string
	Actions    []Action
	TestPaths  []TestPath
}

type TestPath struct {
	Path string
	Name string
}

func createControllerFromDefault(controllerData ControllerData, dest string) error {
	dest = filepath.Join(dest, "controllers")
	lowerName := strings.ToLower(strcase.ToSnake(controllerData.Name))
	controllerPath := filepath.Join(dest, addGoExt(lowerName))
	helpers := []TemplateHelper{
		{
			Name: "whichAction",
			Function: func(action string) string {
				if action == "" {
					log.Println("blank action name provided")
					return ""
				}
				return methodPartial(controllerData, action, "gin")
			},
		},
		{
			Name: "whichActionTest",
			Function: func(action string) string {
				return methodPartial(controllerData, action+"_test", "tests")
			},
		},
	}
	if err := createFileWithHelpers(
		"gin/controller.tmpl", controllerData, controllerPath, helpers); err != nil {
		return err
	}
	// generate controller http tests
	testControllerPath := fmt.Sprintf("%s/%s_test.go", dest, lowerName)
	if err := createFileWithHelpers(
		"tests/controller_test.tpl", controllerData, testControllerPath, helpers); err != nil {
		return err
	}

	// register the controller operations in the router
	routerFilePath := filepath.Join(dest, "router.go")
	AddActionViaAST(controllerData, routerFilePath, dest)

	return nil
}

var methodLookup = map[string]string{
	"GET":    "Show",
	"POST":   "Create",
	"PUT":    "Update",
	"DELETE": "Delete",
}

// OACreateControllerFiles creates a controller for operations found in
// an OpenAPI file
func OACreateControllerFiles(path string, pathItem *openapi3.PathItem, dest string, templateDir string) error {
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
		// i don't remember what this filters
		if op.OperationID == "" && op.Summary == "" {
			log.Printf("No operation ID or summary. Excluding operation from the generated %s controller file\n", name)
			continue
		}
		var handler string
		if isIndex(method, path) {
			handler = "Index"
		} else {
			handler = methodLookup[method]
		}
		action := Action{
			Method: method, Path: path, Handler: handler, Name: handler, Resource: name,
		}
		data.Actions = append(data.Actions, action)
	}

	if err := createControllerFromDefault(data, dest); err != nil {
		return err
	}
	log.Printf("created controller actions for %s\n", path)
	return nil
}
