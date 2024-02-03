package regexp

import (
	"fmt"
	"regexp"
)

const (
	/*
		import quser "exampleproject/internal/app/queries/user" -> user
		%s is 'internal'
	*/
	componentImportRegexpTpl = `".*%s.*\/(.*)"`
)

var (
	/*
		/internal/app/commands/user/create.go -> user
	*/
	componentRegexp = regexp.MustCompile(`.*\/(.*)\/.*\.go`)
)

func CompileComponentImportRegexp(rootNamespace string) (*regexp.Regexp, error) {
	r, err := regexp.Compile(
		fmt.Sprintf(componentImportRegexpTpl, rootNamespace),
	)

	if err != nil {
		return nil, fmt.Errorf("compile: %w", err)
	}

	return r, nil
}

func FindComponent(path string) string {
	matches := componentRegexp.FindStringSubmatch(path)

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
