package component

import (
	"fmt"
	"go/parser"
	"go/token"
)

type Package struct {
	*Component

	imports map[Namespace]*Component
}

func NewPackage(c *Component) *Package {
	return &Package{
		Component: c,
		imports:   make(map[Namespace]*Component),
	}
}

func (p *Package) ParseImportsOfGoFile(goFileName string, componentRegistry *Registry) error {
	file, err := parser.ParseFile(token.NewFileSet(), goFileName, nil, parser.Mode(0))
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	for _, imp := range file.Imports {
		namespace := NewNamespace(imp.Path.Value[1 : len(imp.Path.Value)-1])

		c := componentRegistry.GetOrAddComponent(namespace)

		p.imports[namespace] = c
	}

	return nil
}

func (p *Package) Join(anotherPackage *Package) {
	for namespace, component := range anotherPackage.imports {
		p.imports[namespace] = component
	}
}

func (p *Package) Imports() []*Component {
	var components []*Component

	for _, c := range p.imports {
		components = append(components, c)
	}

	return components
}
