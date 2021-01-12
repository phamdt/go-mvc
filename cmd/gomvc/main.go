package main

import (
	"fmt"
	"os"

	gomvc "github.com/go-generation/go-mvc"

	"github.com/spf13/pflag"
)

// shared flags
var dest string
var spec string
var templateDir string

// db flags
var orm string

func main() {
	root := gomvc.Root()

	app := gomvc.Application()
	root.AddCommand(app)
	setSharedFlags(app.Flags())

	r := gomvc.Resource()
	root.AddCommand(r)
	setSharedFlags(r.Flags())
	r.Flags().StringVarP(&orm, "orm", "o", "", "database access strategy")

	m := gomvc.Model()
	root.AddCommand(m)
	setSharedFlags(m.Flags())
	m.Flags().BoolP("swagger-v2", "2", false, "Swagger v2 spec path")
	m.Flags().BoolP("openapi-v3", "3", false, "OpenAPI v3 spec path")
	m.Flags().StringVarP(&spec, "spec", "s", "./openapi.yml", "OpenAPI spec path")

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
	setSharedFlags(seed.Flags())

	g := gomvc.G()
	root.AddCommand(g)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setSharedFlags(flags *pflag.FlagSet) {
	cwd, _ := os.Getwd()
	flags.StringVarP(&dest, "dest", "d", cwd, "output of generated files")
	flags.StringVarP(&templateDir, "templates", "t", "", "Custom template path")
}
