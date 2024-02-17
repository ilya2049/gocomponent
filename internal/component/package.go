package component

type Package struct {
	*Component

	imports map[Namespace]*Component
}

func NewPackage(c *Component, imports map[Namespace]*Component) *Package {
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

func (p *Package) Imports() []*Component {
	var components []*Component

	for _, component := range p.imports {
		components = append(components, component)
	}

	return components
}
