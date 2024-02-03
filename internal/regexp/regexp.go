package regexp

import (
	"fmt"
	"regexp"
)

const (
	/*
		import "exampleproject/internal/app/queries/user" -> user
		%s is 'internal'
	*/
	singleNameComponentImportRegexpTpl = `".*%s.*\/(.*)"`

	/*
		import "exampleproject/internal/app/queries/user/admin" -> user/admin
		%s is 'internal'
	*/
	doubleNameComponentImportRegexpTpl = `".*%s.*\/(.*\/.*)"`
)

var (
	/*
		/internal/app/commands/user/create.go -> user
	*/
	singleNameComponentRegexp = regexp.MustCompile(`.*\/(.*)\/.*\.go`)

	/*
		/internal/app/commands/user/admin/create.go -> user/admin
	*/
	doubleNameComponentRegexp = regexp.MustCompile(`.*\/(.*\/.*)\/.*\.go`)
)

func CompileComponentImportRegexp(rootNamespace string, isDoubleName bool) (*regexp.Regexp, error) {
	var regexpTpl string

	if isDoubleName {
		regexpTpl = doubleNameComponentImportRegexpTpl
	} else {
		regexpTpl = singleNameComponentImportRegexpTpl
	}

	r, err := regexp.Compile(
		fmt.Sprintf(regexpTpl, rootNamespace),
	)

	if err != nil {
		return nil, fmt.Errorf("compile: %w", err)
	}

	return r, nil
}

func FindComponent(path string, isDoubleName bool) string {
	var matches []string

	if isDoubleName {
		matches = doubleNameComponentRegexp.FindStringSubmatch(path)
	} else {
		matches = singleNameComponentRegexp.FindStringSubmatch(path)
	}

	if len(matches) == 2 {
		return matches[1]
	}

	return ""
}

func FindComponentImports(r *regexp.Regexp, goFileContents string) []string {
	imports := []string{}

	subMatches := r.FindAllStringSubmatch(goFileContents, -1)
	for _, subMatch := range subMatches {
		if len(subMatch) == 2 {
			imports = append(imports, subMatch[1])
		}
	}

	return imports
}
