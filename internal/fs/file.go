package fs

import (
	"path/filepath"
	"regexp"

	"github.com/ilya2049/gocomponent/internal/component"
)

var (
	namespaceRegexp = regexp.MustCompile(`(.*)/.*\.go`)
	goFileRegexp    = regexp.MustCompile(`^.*\.go$`)
)

type fileReader interface {
	ReadFile(name string) ([]byte, error)
}

type filePathWalker interface {
	Walk(root string, fn filepath.WalkFunc) error
}

func findNamespaceInPath(path string) (component.Namespace, bool) {
	matches := namespaceRegexp.FindStringSubmatch(path)

	if len(matches) == 2 {
		return component.NewNamespace(matches[1]), true
	}

	return "", false
}

func isGoSourceFile(path string) bool {
	return goFileRegexp.MatchString(path)
}
