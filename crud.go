package gomvc

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func NewCRUDActions(name string) []Action {
	actions := []Action{}
	for _, action := range []Action{
		{Resource: name, Name: "Index", Method: "GET"},
		{Resource: name, Name: "Create", Method: "POST"},
	} {
		if strings.HasPrefix(name, "/") {
			action.Path = strings.ToLower(name)
		} else {
			action.Path = "/" + strings.ToLower(name)
		}
		action.Handler = strings.Title(action.Name)
		actions = append(actions, action)
	}

	for _, detailAction := range []Action{
		{Resource: name, Name: "Show", Method: "GET"},
		{Resource: name, Name: "Update", Method: "PUT"},
		{Resource: name, Name: "Delete", Method: "DELETE"},
	} {
		detailAction.Path = fmt.Sprintf("/%s/:id", strings.ToLower(name))
		detailAction.Handler = strings.Title(detailAction.Name)
		actions = append(actions, detailAction)
	}
	return actions
}

type Action struct {
	Resource string
	Name     string
	Method   string `json:"method,omitempty"`
	Path     string `json:"path,omitempty"`
	Handler  string `json:"handler,omitempty"`
}

type Response struct {
	Name string
	Code int
	Ref  string
}

// NewResponses creates a list of responses from an OA3 response ref
func NewResponses(specResponses map[string]*openapi3.ResponseRef) []Response {
	var responses []Response
	responseSet := map[string]bool{}
	for statusCode, resRef := range specResponses {
		r := NewResponse(statusCode, resRef)
		if _, ok := responseSet[r.Name]; !ok {
			responseSet[r.Name] = true
			responses = append(responses, r)
		}
	}
	return responses
}

// NewResponse is a constructor for the custom Response object
func NewResponse(statusCode string, resRef *openapi3.ResponseRef) Response {
	code, _ := strconv.Atoi(statusCode)
	return Response{
		Code: code,
		Ref:  resRef.Ref,
		Name: resolveResponseName(resRef),
	}
}

func resolveResponseName(resRef *openapi3.ResponseRef) string {
	if resRef.Ref == "" {
		for _, obj := range resRef.Value.Content {
			name := resolveSchemaName(obj.Schema)
			// TODO: handle multiple
			return name
		}
	}
	return getComponentName(resRef.Ref)
}

func resolveSchemaName(schema *openapi3.SchemaRef) string {
	if schema.Ref == "" {
		return getComponentName(schema.Value.Items.Ref)
	}

	return getComponentName(schema.Ref)
}

func PrintJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "\t")
	log.Println(string(b))
}
