package gomvc

import (
	"fmt"
	"log"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var model = &cobra.Command{
	Use:   "models",
	Short: "Generate model files",
	// Args: func(cmd *cobra.Command, args []string) error {
	// 	if len(args) < 1 {
	// 		return errors.New("requires a name for your new model")
	// 	}
	// 	return nil
	// },
	Run: func(cmd *cobra.Command, args []string) {
		spec, _ := cmd.LocalFlags().GetString("spec")
		if spec == "" {
			panic("you must provide a path to your spec e.g. ./openapi.yml")
		}
		v2, _ := cmd.LocalFlags().GetBool("swagger-v2")
		v3, _ := cmd.LocalFlags().GetBool("openapi-v3")

		var oa3 *openapi3.Swagger
		if v2 && v3 {
			panic("you can only specify one version at a time")
		}
		if v2 {
			oa3 = LoadSwaggerV2AsV3(spec)
		} else {
			// this is defaulted if no version flag is provided
			oa3 = LoadWithKin(spec)
		}
		// TODO: support dest flag
		dir := "./models"
		createDirIfNotExists(dir)
		// create structs in the models dir for all component responses
		for name, r := range oa3.Components.Responses {
			o := LoadResponseObject(name, r.Value)
			str, err := CreateStructFromResponseObject(&o)
			if err != nil {
				panic(err)
			}
			structName := strings.ToLower(strcase.ToSnake(name))
			filename := fmt.Sprintf("%s/%s.go", dir, structName)
			if err := CreateFileFromString(filename, str); err != nil {
				panic(err)
			}
		}

		// create structs in the models dir for all component schemas
		for name, schemaRef := range oa3.Components.Schemas {
			o := LoadSchemaObject(name, schemaRef)
			str, err := CreateStructFromSchemaObject(&o)
			if err != nil {
				panic(err)
			}
			structName := strings.ToLower(strcase.ToSnake(name))
			if err := CreateFileFromString(fmt.Sprintf("%s/%s.go", dir, structName), str); err != nil {
				panic(err)
			}
		}

		// create structs in the models dir for all component schemas
		for name, parameterRef := range oa3.Components.Parameters {
			o := LoadParameterObject(name, parameterRef)
			str, err := CreateStructFromParameterObject(&o)
			if err != nil {
				panic(err)
			}
			structName := strings.ToLower(strcase.ToSnake(name))
			if err := CreateFileFromString(fmt.Sprintf("%s/%s.go", dir, structName), str); err != nil {
				panic(err)
			}
		}

		log.Println("running gofmt and goimports on", dir)
		RunGoFmt(dir)
		RunGoImports(dir)
	},
}

// Model is the cli command that creates new model
func Model() *cobra.Command {
	return model
}
