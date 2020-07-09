package gomvc

import (
	"encoding/json"

	"github.com/aymerick/raymond"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
)

func LoadParameterObject(name string, r *openapi3.ParameterRef) ParameterObject {
	// hack
	b, _ := json.MarshalIndent(r, "", "\t")
	var o ParameterObject
	if err := json.Unmarshal(b, &o); err != nil {
		panic(err)
	}
	SetGoType(&o.Schema)
	return o
}

func pathParamTmpl() string {
	return `
	package models

	type {{camelize Name}} {{Schema.GoType}}
	`
}

func queryParamTmpl() string {
	return `
	package models
	
	type {{camelize Name}} struct
		{{Name}} {{GoType}}
	}
		`
}
func CreateStructFromParameterObject(o *ParameterObject) (string, error) {
	var tmplString string
	if o.In == "path" {
		tmplString = pathParamTmpl()
	} else if o.In == "query" {
		tmplString = queryParamTmpl()
	} else {
		// body
		panic(o.In)
	}
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

type ParameterObject struct {
	Description string   `json:"description,omitempty"`
	In          string   `json:"in,omitempty"`
	Name        string   `json:"name,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Schema      Property `json:"schema,omitempty"`
}

type ParameterSchema struct {
	Format string
	Type   string
	GoType string
}
