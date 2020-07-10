package gomvc

import (
	"encoding/json"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
)

func CreateStructFromResponseObject(m *ResponseModel) (string, error) {
	tmplString := `
package models

type {{Name}} struct {
	{{#if Parent}}
	{{Parent}}{{else}}{{#each Properties}}
	{{#if this.description}}// {{this.description}}{{/if}}
	{{camelize @key}} {{this.GoType}}
	{{/each}}{{/if}}
}
`
	tmpl, err := raymond.Parse(tmplString)
	if err != nil {
		return "", err
	}
	tmpl.RegisterHelper("camelize", func(word string) string {
		return strcase.ToCamel(word)
	})
	result := tmpl.MustExec(m)
	return result, nil
}

func LoadResponseObject(name string, r *openapi3.Response) ResponseModel {
	// hack
	b, _ := json.MarshalIndent(r, "", "\t")
	var o ResponseObject
	if err := json.Unmarshal(b, &o); err != nil {
		panic(err)
	}
	m := ResponseModel{}
	for contentType, content := range o.Content {
		m.Name = name
		m.ContentType = contentType
		// this will be brittel to bad refs
		refParts := strings.Split(content.Schema.Ref, "/")
		refName := refParts[len(refParts)-1]
		if m.Name != refName {
			m.Parent = refName
		}
	}

	return m
}

type ResponseObject struct {
	Content     map[string]SchemaType `json:"content,omitempty"`
	Description string                `json:"description,omitempty"`
}

type ResponseModel struct {
	ContentType string
	Parent      string
	Name        string
	Properties  []Property
}

type SchemaType struct {
	Schema Schema `json:"schema,omitempty"`
}

type Schema struct {
	Ref string `json:"$ref,omitempty"`
}
