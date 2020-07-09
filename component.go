package gomvc

import (
	"encoding/json"
	"fmt"

	"github.com/aymerick/raymond"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
)

func CreateStructFromSchemaObject(o *SchemaObject) (string, error) {
	tmplString := `
package models

type {{Name}} struct
{{#each Properties}}
	{{#if this.description}}// {{this.description}}{{/if}}
	{{camelize @key}} {{this.GoType}}
{{/each}}
}
	`
	tmpl, err := raymond.Parse(tmplString)
	if err != nil {
		return "", err
	}
	tmpl.RegisterHelper("camelize", func(word string) string {
		return strcase.ToCamel(word)
	})

	result := tmpl.MustExec(o)
	return result, nil
}

type Property struct {
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Items       Item   `json:"items,omitempty"`
	Format      string `json:"format,omitempty"`
	Required    bool   `json:"required,omitempty"`
	GoType      string `json:"go_type,omitempty"`
}

type Item struct {
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
}

type SchemaObject struct {
	Name        string
	Description string              `json:"description,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Type        string              `json:"type,omitempty"`
	Title       string              `json:"title,omitempty"`
}

// todo: collect enums
func LoadSchemaObject(name string, r *openapi3.SchemaRef) SchemaObject {
	// hack
	b, _ := json.MarshalIndent(r, "", "\t")
	var o SchemaObject
	if err := json.Unmarshal(b, &o); err != nil {
		panic(err)
	}
	// log.Println(string(b))
	requiredMap := map[string]bool{}
	for _, propertyName := range o.Required {
		requiredMap[propertyName] = true
	}
	for name, property := range o.Properties {
		property.Required = requiredMap[name]
		if property.Format == "" {
			switch property.Type {
			case "boolean":
				property.GoType = "bool"
			case "string":
				property.GoType = "string"
			case "array":
				switch property.Items.Type {
				case "integer":
					property.GoType = "[]int"
				default:
					panic(fmt.Sprintf("%+v", property))
				}
			default:
				panic(fmt.Sprintf("unsupported type: %s", property.Type))
			}
		} else {
			switch property.Format {
			case "date-time":
				property.GoType = "date.Time"
			default:
				property.GoType = property.Format
			}
		}
		o.Properties[name] = property
		if property.GoType == "" {
			panic(property.Type)
		}
	}
	o.Name = strcase.ToCamel(name)
	return o
}
