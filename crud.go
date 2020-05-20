package gomvc

import (
	"fmt"
	"strings"
)

func NewCRUDActions(name string) []Action {
	actions := []Action{}
	for _, action := range []Action{
		{Resource: name, Name: "Index", Method: "GET"},
		{Resource: name, Name: "Create", Method: "POST"},
	} {
		action.Path = fmt.Sprintf("/%s", strings.ToLower(name))
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

type RouteData struct {
	Controllers []ControllerData `json:"controllers,omitempty"`
}
