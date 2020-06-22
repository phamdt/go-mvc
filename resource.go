package gomvc

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

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
		path := filepath.Join("/", strings.ToLower(name))
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
