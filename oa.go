package gomvc

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var oa = &cobra.Command{
	Use:   "oa",
	Short: "Generate controllers from an OpenAPI yml file",
	Run: func(cmd *cobra.Command, args []string) {
		configDir, err := cmd.LocalFlags().GetString("config")
		if err != nil {
			log.Println(err.Error())
			return
		}
		config := NewGoMVCConfig(configDir)

		spec, err := cmd.LocalFlags().GetString("spec")
		if err != nil {
			log.Println(err.Error())
			return
		}
		oa3 := LoadWithKin(spec)

		dest, err := cmd.LocalFlags().GetString("dest")
		if err != nil {
			log.Println(err.Error())
			return
		}
		createDirIfNotExists(dest)
		ctrlDest := filepath.Join(dest, "controllers")
		createDirIfNotExists(ctrlDest)

		templateDir, _ := cmd.LocalFlags().GetString("templates")
		CreateRouter(RouteData{}, "gin/router.tpl", ctrlDest)
		for path, pathItem := range oa3.Paths {
			path = strings.Trim(path, " ")
			log.Printf("examining path, %s\n", path)
			if config.IsBlacklisted(path) {
				continue
			}
			if err := OACreateControllerFiles(path, pathItem, dest, templateDir); err != nil {
				log.Fatalf("%s: %s", path, err.Error())
			}
		}
	},
}

// OA is the cli command that creates a router and controller functions from an
// OpenAPI file
func OA() *cobra.Command {
	return oa
}
