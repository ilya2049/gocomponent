package component

import (
	"fmt"
	"go/parser"
	"go/token"
)

type Package struct {
	*Component

	imports map[string]struct{}
}

func NewPackage(c *Component) *Package {
	return &Package{
		Component: c,
		imports:   make(map[string]struct{}),
	}
}

func (p *Package) ParseImportsOfGoFile(goFileName string) error {
	file, err := parser.ParseFile(token.NewFileSet(), goFileName, nil, parser.Mode(0))
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	for _, imp := range file.Imports {
		p.imports[imp.Path.Value] = struct{}{}
	}

	return nil
}

func (p *Package) Join(another *Package) {
	if !another.namespace.Contains(p.namespace) {
		return
	}

	for imp := range another.imports {
		p.imports[imp] = struct{}{}
	}
}
