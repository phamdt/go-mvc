package main

import (
	"fmt"
	"os"

	gomvc "github.com/go-generation/go-mvc"

	"github.com/spf13/pflag"
)

var dest string
var spec string
var configDir string
var templateDir string

func main() {
	root := gomvc.Root()
	app := gomvc.Application()
	root.AddCommand(app)
	appFlags := app.Flags()
	setSharedFlags(appFlags)

	r := gomvc.Resource()
	root.AddCommand(r)
	resourceFlags := r.Flags()
	setSharedFlags(resourceFlags)

	oa := gomvc.OA()
	root.AddCommand(oa)
	oaFlags := oa.Flags()
	setSharedFlags(oaFlags)
	oaFlags.StringVarP(&spec, "spec", "s", "./openapi.yml", "OpenAPI spec path")

	swagger := gomvc.Swagger()
	root.AddCommand(swagger)
	swaggerFlags := swagger.Flags()
	setSharedFlags(swaggerFlags)
	swaggerFlags.StringVarP(&spec, "spec", "s", "./swagger.yml", "Swagger spec path")

	seed := gomvc.Seed()
	root.AddCommand(seed)
	seedFlags := seed.Flags()
	setSharedFlags(seedFlags)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setSharedFlags(flags *pflag.FlagSet) {
	cwd, _ := os.Getwd()
	flags.StringVarP(&dest, "dest", "d", cwd, "output of generated files")
	flags.StringVarP(&configDir, "config", "c", "", "GoMVC configuration path")
	flags.StringVarP(&templateDir, "templates", "t", "", "Custom template path")
}
