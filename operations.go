package gomvc

import (
	"math"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

func isIndex(method string, path string) bool {
	if strings.ToUpper(method) != "GET" {
		return false
	}

	// TODO: better way to determine?
	return !hasPathParameter(path)
}

func getControllerNameFromPath(path string) string {
	pathParts := strings.Split(path, "/")
	nonParams := []string{}
	for _, part := range pathParts {
		if !hasPathParameter(part) {
			nonParams = append(nonParams, part)
		}
	}
	// limit name to two nouns
	lastTwoIndex := len(nonParams) - 2
	nameIndex := int(math.Max(0, float64(lastTwoIndex)))
	nonParams = nonParams[nameIndex:]
	name := strcase.ToCamel(strings.Join(nonParams, "_"))
	lastPathPart := pathParts[len(pathParts)-1]
	if hasPathParameter(lastPathPart) {
		return inflection.Singular(name)
	}
	return name
}

func hasPathParameter(s string) bool {
	// if it's a param it'll have the format `{paramName}`
	return strings.HasSuffix(s, "}")
}

var methodLookup = map[string]string{
	"GET":    "Show",
	"POST":   "Create",
	"PUT":    "Update",
	"DELETE": "Delete",
}

func getDefaultHandlerName(method, path string) string {
	var handler string
	// index is a special case because currently, all of the HTTP verbs have a 1
	// to 1 relationship with an action except for Index and Show which both use
	// GET.
	if isIndex(method, path) {
		handler = "Index"
	} else {
		handler = methodLookup[method]
	}
	return handler
}
