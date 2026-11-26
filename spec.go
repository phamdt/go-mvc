package gomvc

import (
	"io/ioutil"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	yaml "github.com/ghodss/yaml"
)

// LoadWithKin loads an OpenAPI spec into memory using the kin-openapi library
func LoadWithKin(specPath string) *openapi3.T {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	oa3, err := loader.LoadFromFile(specPath)
	if err != nil {
		panic(err)
	}
	return oa3
}

// LoadSwaggerV2AsV3 takes the file path of a v2 Swagger file and returns a
// V3 representation
func LoadSwaggerV2AsV3(specPath string) *openapi3.T {
	swaggerSpec := openapi2.T{}
	c, err := ioutil.ReadFile(specPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(c, &swaggerSpec)
	if err != nil {
		panic(err)
	}
	oa3, err := openapi2conv.ToV3(&swaggerSpec)
	if err != nil {
		panic(err)
	}
	return oa3
}

// assumes usage of OpenAPI 3.x spec in which component refs are formatted as
// '#/componeents/<sub component type>/<user defined component name>'
func getComponentName(ref string) string {
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}
