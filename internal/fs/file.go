package fs

import (
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
