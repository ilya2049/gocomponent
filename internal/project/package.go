package project

import "github.com/ilya2049/gocomponent/internal/component"

type Package struct {
	*component.Component

	imports map[component.Namespace]*component.Component
}

func NewPackage(c *component.Component, imports map[component.Namespace]*component.Component) *Package {
	return &Package{
		Component: c,
		imports:   imports,
	}
}

func (p *Package) Join(anotherPackage *Package) {
	for namespace, component := range anotherPackage.imports {
		p.imports[namespace] = component
	}
}

func (p *Package) Imports() []*component.Component {
	var components []*component.Component

	for _, component := range p.imports {
		components = append(components, component)
	}

	return components
}
