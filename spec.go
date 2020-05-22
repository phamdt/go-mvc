package gomvc

import "github.com/getkin/kin-openapi/openapi3"

// LoadWithKin loads an OpenAPI spec into memory using the kin-openapi library
func LoadWithKin(specPath string) *openapi3.Swagger {
	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	oa3, err := loader.LoadSwaggerFromFile(specPath)
	if err != nil {
		panic(err)
	}
	return oa3
}
