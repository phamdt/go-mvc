package gomvc

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
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
		template := args[0]
		name := args[1]
		log.Printf("preparing to create a new %s: %s\n", template, name)

		baseTemplateDir := "./templates"
		templateDir := fmt.Sprintf("%s/%s.tpl", baseTemplateDir, template)
		destPath := fmt.Sprintf("%s/%s.%s.go", ".", name, template)
		packageName := GetPackageName()
		data := map[string]string{
			"Package":   packageName,
			"Name":      name,
			"TitleName": strings.Title(name),
		}

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
	},
}

func GetPackageName() string {
	cwd, _ := os.Getwd()
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
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
