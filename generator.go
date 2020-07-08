package gomvc

import (
	"log"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

// Generator wraps functionality for reading and manipulating a single OpenAPI spec
type Generator struct {
	oa3 *openapi3.Swagger
}

// NewGenerator is a constructor for Generator
func NewGenerator(oa3 *openapi3.Swagger) Generator {
	return Generator{
		oa3: oa3,
	}
}

// CreateControllerFiles creates a controller for operations found in
// an OpenAPI file
func (oag *Generator) CreateControllerFiles(path string, pathItem *openapi3.PathItem, dest string, templateDir string) error {
	name := strcase.ToSnake(pathItem.Summary)
	name = strings.ToLower(name)

	if name == "" {
		log.Printf("No summary provided in API, defaulting to deriving name from path %s since we can't identify a name for the resource", path)
		name = getControllerNameFromPath(path)
	}

	log.Printf("Preparing to generate controller files for %s %s", name, path)
	data := ControllerData{
		Name:       name,
		PluralName: inflection.Plural(name),
		Path:       path,
		Actions:    []Action{},
	}

	var responses []Response
	responseSet := map[string]bool{}

	// collect controller methods based on specified HTTP verbs/operations
	for method, op := range pathItem.Operations() {
		var handler = getDefaultHandlerName(method, path)
		var operationName string
		if op.OperationID == "" {
			log.Printf("Missing operation ID. Generating default name for handler/operation function in controller %s.\n", name)
			operationName = handler
		} else {
			operationName = strings.Title(op.OperationID)
		}
		opResponses := NewResponses(op.Responses)

		// responses might not be unique across controller actions
		// so this ensures that any associated generation happens once
		// additionally, we only care about the error responses
		// to generate filter functions
		for _, r := range opResponses {
			if r.Code < 400 {
				continue
			}
			if _, ok := responseSet[r.Name]; !ok {
				responseSet[r.Name] = true
				responses = append(responses, r)
			}
		}

		a := Action{
			SingularResource: inflection.Singular(name),
			Resource:         name,
			Method:           method,
			Path:             path,
			Handler:          handler,
			Name:             operationName,
		}
		data.Actions = append(data.Actions, a)
	}
	data.ErrorResponses = responses

	if err := createControllerFromDefault(data, dest); err != nil {
		return err
	}
	log.Printf("created controller actions for %s\n", path)
	return nil
}
