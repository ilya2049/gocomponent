package fs

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ilya2049/gocomponent/internal/domain/component"
)

var (
	namespaceRegexp = regexp.MustCompile(`(.*)/.*\.go`)
	goFileRegexp    = regexp.MustCompile(`^.*\.go$`)
)

type filePathWalker interface {
	Walk(root string, fn filepath.WalkFunc) error
}

type astFileParser interface {
	ParseFile(fset *token.FileSet, filename string, src any, mode parser.Mode) (f *ast.File, err error)
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

func isHidden(path string) bool {
	return strings.HasPrefix(path, ".")
}
