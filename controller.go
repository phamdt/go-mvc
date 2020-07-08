package gomvc

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

type ControllerData struct {
	Name           string
	PluralName     string
	Path           string
	Actions        []Action
	TestPaths      []TestPath
	ErrorResponses []Response
	ORM            string
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
			Function: func(handler string) string {
				if handler == "" {
					log.Println("blank handler name provided")
					return ""
				}
				actionData := findActionByHandler(controllerData.Actions, handler)
				return methodPartial(actionData, handler, "gin")
			},
		},
		{
			Name: "whichActionTest",
			Function: func(handler string) string {
				actionData := findActionByHandler(controllerData.Actions, handler)
				return methodPartial(actionData, handler+"_test", "tests")
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

// find specific action tied to the handler
func findActionByHandler(actions []Action, handler string) Action {
	var current Action
	for _, a := range actions {
		if a.Handler == handler {
			current = a
			break
		}
	}
	return current
}
