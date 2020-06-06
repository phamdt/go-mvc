package gomvc

import (
	"log"

	"github.com/spf13/cobra"
)

var swagger = &cobra.Command{
	Use:   "swagger",
	Short: "Generate controllers from a v2 Swagger yml file",
	Run: func(cmd *cobra.Command, args []string) {
		configDir, err := cmd.LocalFlags().GetString("config")
		if err != nil {
			log.Println(err.Error())
			return
		}
		// TODO: read spec location from config
		specPath, err := cmd.LocalFlags().GetString("spec")
		if err != nil {
			log.Println(err.Error())
			return
		}

		// read intended destination for generation output
		destDir, err := cmd.LocalFlags().GetString("dest")
		if err != nil {
			log.Println(err.Error())
			return
		}

		templateDir, err := cmd.LocalFlags().GetString("templates")
		if err != nil {
			log.Println(err.Error())
			return
		}
		oa3 := LoadSwaggerV2AsV3(specPath)
		GenerateFromOA(oa3, destDir, templateDir, configDir)
	},
}

// Swagger is the cli command that creates a router and controller functions
// from a swagger file
func Swagger() *cobra.Command {
	return swagger
}
