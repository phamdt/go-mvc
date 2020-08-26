package gomvc

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

var oa = &cobra.Command{
	Use:   "oa",
	Short: "Generate controllers from an OpenAPI yml file",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.LocalFlags()
		configDir, err := flags.GetString("config")
		if err != nil {
			log.Println(err.Error())
			return
		}
		// TODO: read spec location from config
		specPath, err := flags.GetString("spec")
		if err != nil {
			log.Println(err.Error())
			return
		}

		// read intended destination for generation output
		destDir, err := flags.GetString("dest")
		if err != nil {
			log.Println(err.Error())
			return
		}

		templateDir, err := flags.GetString("templates")
		if err != nil {
			log.Println(err.Error())
			return
		}
		oa3 := LoadWithKin(specPath)
		GenerateFromOA(oa3, destDir, templateDir, configDir)
	},
}

// GenerateFromOA is the primary logic for the oa command, creating controllers
func GenerateFromOA(oa3 *openapi3.Swagger, dest, templateDir, configDir string) {
	config := NewGoMVCConfig(configDir)

	createDirIfNotExists(dest)
	ctrlDest := filepath.Join(dest, "controllers")
	createDirIfNotExists(ctrlDest)

	CreateRouter(RouteData{}, "gin/router.tpl", ctrlDest)
	g := NewGenerator(oa3)
	for path, pathItem := range oa3.Paths {
		path = strings.Trim(path, " ")
		log.Printf("examining path, %s\n", path)
		if config.IsDenylisted(path) {
			continue
		}
		if err := g.CreateControllerFiles(path, pathItem, dest, templateDir); err != nil {
			log.Fatalf("%s: %s", path, err.Error())
		}
	}
}

// OA is the cli command that creates a router and controller functions from an
// OpenAPI file
func OA() *cobra.Command {
	return oa
}
