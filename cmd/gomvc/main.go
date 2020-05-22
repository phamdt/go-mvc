package main

import (
	"fmt"
	"os"

	gomvc "github.com/phamdt/go-mvc"

	"github.com/spf13/pflag"
)

var dest string
var spec string
var configDir string
var templateDir string

func main() {
	root := gomvc.Root
	root.AddCommand(gomvc.Application)
	appFlags := gomvc.Application.Flags()
	setSharedFlags(appFlags)

	root.AddCommand(gomvc.Resource)
	resourceFlags := gomvc.Resource.Flags()
	setSharedFlags(resourceFlags)

	root.AddCommand(gomvc.OA)
	oaFlags := gomvc.OA.Flags()
	setSharedFlags(oaFlags)
	oaFlags.StringVarP(&spec, "spec", "s", "./openapi.yml", "OpenAPI spec path")

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setSharedFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&dest, "dest", "d", "", "output of generated files")
	flags.StringVarP(&configDir, "config", "c", "", "GoMVC configuration path")
	flags.StringVarP(&templateDir, "templates", "t", "", "Custom template path")
}
