package gomvc

import (
	"fmt"
	"log"
	"strings"

	"github.com/hoisie/mustache"
)

func NewCRUDActions(name string) []Action {
	actions := []Action{}
	for _, action := range []Action{
		Action{Resource: name, Name: "Index", Method: "GET"},
		Action{Resource: name, Name: "Create", Method: "POST"},
	} {
		action.Path = fmt.Sprintf("/%s", strings.ToLower(name))
		action.Handler = strings.Title(action.Name)
		actions = append(actions, action)
	}

	for _, detailAction := range []Action{
		Action{Resource: name, Name: "Show", Method: "GET"},
		Action{Resource: name, Name: "Update", Method: "PUT"},
		Action{Resource: name, Name: "Delete", Method: "DELETE"},
	} {
		detailAction.Path = fmt.Sprintf("/%s/:id", strings.ToLower(name))
		detailAction.Handler = strings.Title(detailAction.Name)
		actions = append(actions, detailAction)
	}
	return actions
}

type Action struct {
	Resource string
	Name     string
	Method   string `json:"method,omitempty"`
	Path     string `json:"path,omitempty"`
	Handler  string `json:"handler,omitempty"`
}

type RouteData struct {
	Controllers []ControllerData `json:"controllers,omitempty"`
}

type ControllerData struct {
	Name       string
	PluralName string
	Path       string
	Actions    []Action
}

// CreateControllerFrom actually creates files with controllerData presumably parsed, processed from some source of truth regarding API design e.g. OpenAPI or a CLI wizard
func CreateControllerFrom(controllerData ControllerData, templateDir string, dest string) error {
	// generate controller with controllerData
	controllerDir := fmt.Sprintf("%s/gin/controller.tmpl", templateDir)
	r := mustache.RenderFile(controllerDir, controllerData)
	filepath := fmt.Sprintf("%s/%s.go", dest, controllerData.Name)
	err := createFileFromString(filepath, r)
	if err != nil {
		log.Println("error generating file for", filepath, err.Error())
		return err
	}
	// generate controller http tests
	testPath := fmt.Sprintf("%s/tests/controller_test.tpl", templateDir)
	tr := mustache.RenderFile(testPath, controllerData)
	testfilepath := fmt.Sprintf("%s/%s_test.go", dest, controllerData.Name)
	if err := createFileFromString(testfilepath, tr); err != nil {
		log.Println("error generating test file for", testfilepath, err.Error())
		return err
	}
	return nil
}
