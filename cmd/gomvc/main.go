package main

import (
	"fmt"
	gomvc "github.com/phamdt/go-mvc"
	"os"

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

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setSharedFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&dest, "dest", "d", "", "output of generated files")
	// flags.StringVarP(&spec, "spec", "s", "./openapi.yml", "OpenAPI spec path")
	// flags.StringVarP(&configDir, "config", "c", "", "GoMVC configuration path")
	// flags.StringVarP(&templateDir, "template-dir", "t", "", "Custom template path")
}
