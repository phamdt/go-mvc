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
