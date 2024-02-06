package impset

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type NamespaceImportSet struct {
	namespace string
	imports   map[string]struct{}
}

func NewNamespaceImportSet(namespace string) *NamespaceImportSet {
	return &NamespaceImportSet{
		namespace: namespace,
		imports:   make(map[string]struct{}),
	}
}

func (s *NamespaceImportSet) ParseImportsOfGoFile(goFileName string) error {
	file, err := parser.ParseFile(token.NewFileSet(), goFileName, nil, parser.Mode(0))
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	for _, imp := range file.Imports {
		s.imports[imp.Path.Value] = struct{}{}
	}

	return nil
}

func (s *NamespaceImportSet) Join(another *NamespaceImportSet) {
	if !strings.Contains(another.namespace, s.namespace) {
		return
	}

	for imp := range another.imports {
		s.imports[imp] = struct{}{}
	}
}
