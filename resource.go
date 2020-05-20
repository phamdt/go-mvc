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
