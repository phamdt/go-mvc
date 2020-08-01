package gomvc

import (
	"fmt"
	"log"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/aymerick/raymond"
)

type TemplateHelper struct {
	Name     string
	Function func(string) string
}

// TODO: consolidate this with createFileFromTemplates
func createContentsFromTemplate(tmplPath string, data interface{}) string {
	tmpl, err := getTemplate(tmplPath)
	if err != nil {
		panic(err)
	}
	result := tmpl.MustExec(data)
	return result
}

// TODO there's some duplication here and in the helper registration in controller.go
func createFileFromTemplates(template string, data interface{}, destPath string) error {
	content := createContentsFromTemplate(template, data)
	if err := CreateFileFromString(destPath, content); err != nil {
		log.Println("could not create file for", destPath)
		return err
	}
	return nil
}

// TODO: support custom templates
func methodPartial(ctx interface{}, name string, subDir string) string {
	name = strings.ToLower(name)
	tmplPath := fmt.Sprintf("%s/partials/%s.tmpl", subDir, name)

	tmpl, err := getTemplate(tmplPath)
	if err != nil {
		panic(err)
	}
	return tmpl.MustExec(ctx)
}

func createFileWithHelpers(tmplPath string, data interface{}, destPath string, helpers []TemplateHelper) error {
	tmpl, err := getTemplate(tmplPath)
	if err != nil {
		return err
	}

	for _, helper := range helpers {
		tmpl.RegisterHelper(helper.Name, helper.Function)
	}
	interpolated := tmpl.MustExec(data)
	if err := CreateFileFromString(destPath, interpolated); err != nil {
		log.Println("could not create file for", destPath)
		return err
	}
	return nil
}

func getTemplate(tmplPath string) (*raymond.Template, error) {
	box := rice.MustFindBox("templates")
	tmplString := box.MustString(tmplPath)
	return raymond.Parse(tmplString)
}
