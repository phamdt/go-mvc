package gomvc

import (
	"math"
	"strings"

	"github.com/iancoleman/strcase"
)

func isIndex(method string, path string) bool {
	if strings.ToUpper(method) != "GET" {
		return false
	}

	// TODO: better way to determine?
	return !strings.HasSuffix(path, "}")
}

func getControllerNameFromPath(path string) string {
	pathParts := strings.Split(path, "/")
	nonParams := []string{}
	for _, part := range pathParts {
		// if it's a param it'll have the format `{paramName}`
		if !strings.HasSuffix(part, "}") {
			nonParams = append(nonParams, part)
		}
	}
	lastTwoIndex := len(nonParams) - 2
	nameIndex := int(math.Max(0, float64(lastTwoIndex)))
	nonParams = nonParams[nameIndex:]
	name := strcase.ToCamel(strings.Join(nonParams, "_"))
	return name
}
