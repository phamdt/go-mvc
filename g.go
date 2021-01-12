package gomvc

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/spf13/cobra"
)

var g = &cobra.Command{
	Use:   "g",
	Short: "Generate files from arbitrary templates",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			msg := "usage: gomvc g <your template name> <object name>"
			return errors.New(msg)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		templateType := args[0]
		name := args[1]
		log.Printf("preparing to create %s file(s) for %s\n", templateType, name)
		if err := GenerateFromUserTemplate(name, templateType); err != nil {
			panic(err)
		}
	},
}

func CreateFileFromLocalTemplate(data interface{}, templateDir, destPath string) {
	t, err := raymond.ParseFile(templateDir)
	if err != nil {
		panic(err)
	}
	r, err := t.Exec(data)
	if err != nil {
		panic(err)
	}
	if err := CreateFileFromString(destPath, r); err != nil {
		panic(err)
	}
}

func GenerateFromUserTemplate(name string, templateType string) error {
	baseTemplateDir := "./templates"
	src := fmt.Sprintf("%s/%s", baseTemplateDir, templateType)
	packageName := GetPackageName()
	data := map[string]string{
		"Package":   packageName,
		"Name":      name,
		"TitleName": strings.Title(name),
	}

	if dirExists(src) {
		if err := handleDir(src, data); err != nil {
			return err
		}
	} else if fileExists(src + ".tpl") {
		templateDir := fmt.Sprintf("%s/%s.tpl", baseTemplateDir, templateType)
		destPath := fmt.Sprintf("%s/%s.go", ".", name)

		CreateFileFromLocalTemplate(data, templateDir, destPath)
	} else {
		log.Printf("'%s' is neither a template or directory of templates\n", src)
	}
	return nil
}

func handleDir(src string, data interface{}) error {
	// create a file from all templates in the directory recursively
	templates, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, f := range templates {
		if f.IsDir() {
			if err := handleDir(filepath.Join(src, f.Name()), data); err != nil {
				return err
			}
		}
		templateDir := filepath.Join(src, f.Name())
		parts := strings.Split(f.Name(), ".")
		nameWithoutExt := parts[len(parts)-2]
		destPath := fmt.Sprintf("%s/%s.go", ".", nameWithoutExt)

		CreateFileFromLocalTemplate(data, templateDir, destPath)
	}
	return nil
}

func GetPackageName() string {
	cwd, _ := os.Getwd()
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, cwd, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parse dir error: %v\n ", err)
	}
	for _, pkg := range pkgs {
		return pkg.Name
	}
	return ""
}

// G is the cli command that creates new g
func G() *cobra.Command {
	return g
}
