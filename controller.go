package gomvc

import (
	"fmt"
	"log"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/aymerick/raymond"
	"github.com/davecgh/go-spew/spew"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type ControllerData struct {
	Name       string
	PluralName string
	Path       string
	Actions    []Action
}

func createControllerFromDefault(controllerData ControllerData, dest string) error {
	box := rice.MustFindBox("templates")
	tmplString, err := box.String("gin/controller.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := raymond.Parse(tmplString)
	if err != nil {
		return err
	}
	testTmplString := box.MustString("tests/controller_test.tpl")
	testTmpl, err := raymond.Parse(testTmplString)
	if err != nil {
		return err
	}
	raymond.RegisterHelper("whichAction", func(action string) string {
		log.Println("looking for HTTP action partial", action)
		if action == "" {
			log.Println("blank action name provided")
			return ""
		}
		return methodPartial(controllerData, action, "gin")
	})
	lowerName := strings.ToLower(controllerData.Name)
	filepath := fmt.Sprintf("%s/%s.go", dest, lowerName)
	log.Println("fp", filepath)
	result := tmpl.MustExec(controllerData)
	if err := createFileFromString(filepath, result); err != nil {
		log.Println("error generating file for", filepath, err.Error())
		return err
	}
	// register the controller operations in the router
	routerFilePath := fmt.Sprintf("%s/router.go", dest)
	AddActionViaAST(controllerData.Actions, routerFilePath, dest)

	// generate controller http tests
	raymond.RegisterHelper("whichActionTest", func(action string) string {
		log.Println("looking for HTTP action test partial", action)
		return methodPartial(controllerData, action+"_test", "tests")
	})
	testfilepath := fmt.Sprintf("%s/%s_test.go", dest, lowerName)
	testResult := testTmpl.MustExec(controllerData)
	if err := createFileFromString(testfilepath, testResult); err != nil {
		log.Println("error generating file for", testfilepath, err.Error())
		return err
	}
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
	// skip creating file if we can't find a good name from the doc
	if name == "" {
		log.Printf("No summary provided in API, defaulting to deriving name from path %s since we can't identify a name for the resource", path)
		name = getControllerNameFromPath(path)
	}
	log.Printf("Preparing to generate controller files for %s", name)

	data := ControllerData{
		Name:       inflection.Singular(name),
		PluralName: inflection.Plural(name),
		Path:       path,
		Actions:    []Action{},
	}

	// collect controller methods based on specified HTTP verbs/operations
	for method, op := range pathItem.Operations() {
		// i don't remember what this filters
		if op.OperationID == "" && op.Summary == "" {
			log.Printf("no operation ID or summary. Excluding operation from the generated %s controller file\n", name)
			continue
		}
		var handler string
		if isIndex(method, path) {
			handler = "Index"
		} else {
			handler = methodLookup[method]
		}
		action := Action{
			Method: method, Path: path, Handler: handler, Name: handler,
		}
		data.Actions = append(data.Actions, action)
		spew.Dump(data)
	}

	err := createControllerFromDefault(data, dest)
	if err != nil {
		return err
	}

	log.Printf("created %s\n", name)
	return nil
}
