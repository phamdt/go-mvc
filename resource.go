package gomvc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/aymerick/raymond"
	"github.com/jinzhu/inflection"
	"github.com/spf13/cobra"
)

var resource = &cobra.Command{
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
		path := filepath.Join(name)
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
			dest = path
		}

		if err := createControllerFromDefault(controllerData, dest); err != nil {
			panic(err)
		}
	},
}

// Resource is the cli command that creates new resource
func Resource() *cobra.Command {
	return resource
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
