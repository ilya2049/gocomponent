package fswalk

import (
	"regexp"
)

const (
	// import quser "exampleproject/internal/app/queries/user" -> user
	// %s is 'internal'
	componentConnectionRegexpTemplate = `".*%s.*\/(.*)"`
)

var (
	// /internal/app/commands/user/create.go -> user
	componentRegexp = regexp.MustCompile(`.*\/(.*)\/.*\.go`)
)

func getComponentUsingRegexp(path string) string {
	regexpMatches := componentRegexp.FindStringSubmatch(path)

	if len(regexpMatches) == 2 {
		return regexpMatches[1]
	}

	return ""
}

func getComponentConnectionsUsingRegexp(
	componentConnectionRegexp *regexp.Regexp,
	goFileContents string,
) []string {
	componentConnections := []string{}

	subMatches := componentConnectionRegexp.FindAllStringSubmatch(goFileContents, -1)
	for _, subMatch := range subMatches {
		if len(subMatch) == 2 {
			componentConnections = append(componentConnections, subMatch[1])
		}
	}

	return componentConnections
}
