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
