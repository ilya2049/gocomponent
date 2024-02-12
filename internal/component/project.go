package component

type Project struct {
	packages map[Namespace]*Package
}

func NewProject() *Project {
	return &Project{
		packages: make(map[Namespace]*Package),
	}
}

func (p *Project) FindPackage(namespace Namespace) (*Package, bool) {
	pkg, ok := p.packages[namespace]

	return pkg, ok
}

func (p *Project) AddPackage(namespace Namespace, pkg *Package) {
	p.packages[namespace] = pkg
}

func (p *Project) Packages() []*Package {
	packages := make([]*Package, 0, len(p.packages))

	for _, p := range p.packages {
		packages = append(packages, p)
	}

	return packages
}

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
