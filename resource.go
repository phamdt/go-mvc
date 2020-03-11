package gomvc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/aymerick/raymond"
	"github.com/jinzhu/inflection"
	"github.com/spf13/cobra"
)

// Resource is the cli command that creates new resource
var Resource = &cobra.Command{
	Use:   "resource",
	Short: "Generate resource files",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name for your new resource")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		log.Printf("preparing to create a new resource %s\n", name)
		path := fmt.Sprintf("/%s", name)
		controllerData := ControllerData{
			Name:       strings.Title(inflection.Singular(name)),
			PluralName: inflection.Plural(name),
			Path:       path,
			Actions:    NewCRUDActions(name),
		}
		dest, _ := cmd.LocalFlags().GetString("dest")
		if dest == "" {
			path, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			dest = fmt.Sprintf("%s/controllers", path)
		}

		createControllerFromDefault(controllerData, dest)
	},
}

// TODO: support custom templates
func methodPartial(ctx interface{}, name string, subDir string) string {
	name = strings.ToLower(name)
	box := rice.MustFindBox("templates")
	tmplDir := fmt.Sprintf("%s/partials/%s.tmpl", subDir, name)
	t := box.MustString(tmplDir)
	tmpl, err := raymond.Parse(t)
	if err != nil {
		panic(err)
	}
	return tmpl.MustExec(ctx)
}

func createControllerFromDefault(controllerData ControllerData, dest string) {
	box := rice.MustFindBox("templates")
	tmplString := box.MustString("gin/controller.tmpl")
	tmpl, err := raymond.Parse(tmplString)
	if err != nil {
		panic(err)
	}
	testTmplString := box.MustString("tests/controller_test.tpl")
	testTmpl, err := raymond.Parse(testTmplString)
	if err != nil {
		panic(err)
	}
	baseCreateController(controllerData, dest, tmpl, testTmpl)
}

func createControllerFromHandleBars(controllerData ControllerData, templateDir string, dest string) error {
	// generate controller with controllerData
	controllerDir := fmt.Sprintf("%s/gin/controller.tmpl", templateDir)
	tmpl, err := raymond.ParseFile(controllerDir)
	if err != nil {
		panic(err)
	}
	controllerTestDir := fmt.Sprintf("%s/tests/controller_test.tpl", templateDir)
	testTmpl, err := raymond.ParseFile(controllerTestDir)
	if err != nil {
		panic(err)
	}
	return baseCreateController(controllerData, dest, tmpl, testTmpl)
}

func baseCreateController(controllerData ControllerData, dest string, ctrlTpl, testTpl *raymond.Template) error {
	raymond.RegisterHelper("whichAction", func(action string) string {
		log.Println("looking for HTTP action partial", action)
		return methodPartial(controllerData, action, "gin")
	})
	filepath := fmt.Sprintf("%s/%s.go", dest, strings.ToLower(controllerData.Name))
	result := ctrlTpl.MustExec(controllerData)
	if err := createFileFromString(filepath, result); err != nil {
		log.Println("error generating file for", filepath, err.Error())
		return err
	}

	routerFilePath := fmt.Sprintf("%s/router.go", dest)
	AddActionViaAST(controllerData.Actions, routerFilePath, dest)

	// generate controller http tests
	raymond.RegisterHelper("whichActionTest", func(action string) string {
		log.Println("looking for HTTP action test partial", action)
		return methodPartial(controllerData, action+"_test", "tests")
	})
	testfilepath := fmt.Sprintf("%s/%s_test.go", dest, strings.ToLower(controllerData.Name))
	testResult := testTpl.MustExec(controllerData)
	if err := createFileFromString(testfilepath, testResult); err != nil {
		log.Println("error generating file for", testfilepath, err.Error())
		return err
	}
	return nil
}
