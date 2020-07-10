package gomvc

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
)

func CreateStructFromSchemaObject(o *SchemaObject) (string, error) {
	tmplString := `
package models

{{#if GoType}}
type {{Name}} {{GoType}}
{{else}}
type {{Name}} struct {
{{#each Properties}}
{{#if this.description}}
// {{this.description}}{{/if}}
	{{camelize @key}} {{this.GoType}}{{/each}}
}
{{/if}}
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
	Name                 string               `json:"name,omitempty"`
	Description          string               `json:"description,omitempty"`
	Type                 string               `json:"type,omitempty"`
	Items                Item                 `json:"items,omitempty"`
	Format               string               `json:"format,omitempty"`
	Required             bool                 `json:"required,omitempty"`
	GoType               string               `json:"go_type,omitempty"`
	Ref                  string               `json:"$ref,omitempty"`
	AdditionalProperties AdditionalProperties `json:"additionalProperties,omitempty"`
}

type AdditionalProperties struct {
	Type  string
	OneOf []Type
	AnyOf []Type
}

type Type struct {
	Type string
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
	GoType      string
	Items       Item
}

// todo: collect enums
func LoadSchemaObject(name string, r *openapi3.SchemaRef) SchemaObject {
	// hack
	b, _ := json.MarshalIndent(r, "", "\t")
	var o SchemaObject
	if err := json.Unmarshal(b, &o); err != nil {
		panic(err)
	}
	if o.Type == "array" {
		o.GoType = fmt.Sprintf("[]%s", o.Items.Type)
	} else if objectType(o.Type) {
		requiredMap := map[string]bool{}
		for _, propertyName := range o.Required {
			requiredMap[propertyName] = true
		}
		for name, property := range o.Properties {
			property.Required = requiredMap[name]
			// mutation
			goType := GetGoType(name, &property)
			property.GoType = goType
			o.Properties[name] = property
		}
	} else {
		o.GoType = GetGoPrimitiveType(o.Type)
	}
	o.Name = strcase.ToCamel(name)
	return o
}

func objectType(t string) bool {
	return t == "object"
}

func GetGoPrimitiveType(t string) string {
	switch t {
	case "boolean":
		return "bool"
	case "string":
		return "string"
	case "integer":
		return "int"
	default:
		panic(fmt.Sprintf("invalid primitive: %s", t))
	}
}

// GetGoType mutates the property type to determine GoType
func GetGoType(name string, p *Property) string {
	if p.Format == "" {
		switch p.Type {
		case "boolean":
			return "bool"
		case "string":
			return "string"
		case "integer":
			return "int"
		case "array":
			switch p.Items.Type {
			case "integer":
				return "[]int"
			case "number", "double":
				return "float"
			default:
				panic(fmt.Sprintf("unexpected format: %s", p.Format))
			}
		case "object":
			if hasAdditionalProperties(p) {
				return "map[string]interface{}"
			}
			refObjectName := GetLastPathPart(p.Ref)
			return refObjectName
		case "number", "double":
			return "float64"
		case "":
			if p.Ref == "" {
				panic(fmt.Sprintf("while your spec may be valid, we require properties to have a type (e.g. integer, array) for %s %s", name, p.Ref))
			}
			// is this enough? what if there's a ref and also properties
			refObjectName := GetLastPathPart(p.Ref)
			return refObjectName
		default:
			panic(fmt.Sprintf("unsupported type: %s for %s", p.Type, name))
		}
	} else {
		switch p.Format {
		case "date-time":
			return "time.Time"
		case "number", "double":
			return "float64"
		default:
			return p.Format
		}
	}
}

func GetLastPathPart(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func hasAdditionalProperties(p *Property) bool {
	ap := &p.AdditionalProperties
	return ap.Type != "" && len(ap.AnyOf) == 0 && len(ap.OneOf) == 0
}
