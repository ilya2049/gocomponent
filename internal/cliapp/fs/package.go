package fs

import "github.com/ilya2049/gocomponent/internal/component"

type projectPackage struct {
	*component.Component

	imports map[component.Namespace]*component.Component
}

func newProjectPackage(c *component.Component, imports map[component.Namespace]*component.Component) *projectPackage {
	return &projectPackage{
		Component: c,
		imports:   imports,
	}
}

func (p *projectPackage) join(anotherPackage *projectPackage) {
	for namespace, component := range anotherPackage.imports {
		p.imports[namespace] = component
	}
}

func (p *projectPackage) getImports() []*component.Component {
	var components []*component.Component

	for _, component := range p.imports {
		components = append(components, component)
	}

	return components
}
