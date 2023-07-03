package router

import (
	"goHtmlBuilder/filescaner"
	"strings"
)

func BuildRoutes(ghtmlFiles []string) (map[string]string, error) {
	result := make(map[string]string)

	for _, path := range ghtmlFiles {
		result[strings.Replace(path, filescaner.GhtmlExt, "", -1)] = "/dist/" + strings.Replace(path, filescaner.GhtmlExt, ".html", -1)
	}
	return result, nil
}
