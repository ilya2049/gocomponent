package fs

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ilya2049/gocomponent/internal/component"
)

var namespaceRegexp = regexp.MustCompile(`(.*)/.*\.go`)

func findNamespaceInPath(path string) (component.Namespace, bool) {
	matches := namespaceRegexp.FindStringSubmatch(path)

	if len(matches) == 2 {
		return component.NewNamespace(matches[1]), true
	}

	return "", false
}

func isGoSourceFile(path string) bool {
	return strings.Contains(path, ".go")
}

var ErrFirstLineOfGoModShouldIncludeExactlyTwoPArts = errors.New(
	"the first line of go mod file parts should includes exactly two parts",
)

func readModuleName(projectDir string) (string, error) {
	goModFileContents, err := os.ReadFile(projectDir + "go.mod")
	if err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}

	var firstLineOfGoModFile = []byte{}

	for _, b := range goModFileContents {
		if b == '\n' {
			break
		} else {
			firstLineOfGoModFile = append(firstLineOfGoModFile, b)
		}
	}

	firstLineOfGoModFileParts := strings.Split(string(firstLineOfGoModFile), " ")
	if len(firstLineOfGoModFileParts) != 2 {
		return "", ErrFirstLineOfGoModShouldIncludeExactlyTwoPArts
	}

	return firstLineOfGoModFileParts[1], nil
}
