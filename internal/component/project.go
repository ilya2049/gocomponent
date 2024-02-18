package component

type Project struct {
	packages map[Namespace]*Package

	components map[Namespace]*Component
}

func NewProject() *Project {
	return &Project{
		packages:   make(map[Namespace]*Package),
		components: make(map[Namespace]*Component),
	}
}

func (p *Project) GetOrAddComponent(namespace Namespace) *Component {
	existingComponent, ok := p.components[namespace]
	if ok {
		return existingComponent
	}

	newComponent := New(namespace)

	p.components[namespace] = newComponent

	return newComponent
}

func (p *Project) Components() Components {
	var components Components

	for _, component := range p.components {
		components = append(components, component)
	}

	return components
}

func (p *Project) MakeUniqueComponentIDs() {
	components := p.Components()

	for len(components) > 0 {
		firstComponent := components[0]
		isComponentIDUnique := true

		for i := 1; i < len(components); i++ {
			if components[i].ID() == firstComponent.ID() {
				isComponentIDUnique = false
				components[i].ExtendID()
			}
		}

		if isComponentIDUnique {
			components = components[1:]
		} else {
			firstComponent.ExtendID()
		}
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

func (p *Project) CreateComponentGraph() *Graph {
	imports := make(Imports, 0)

	for _, p := range p.packages {
		for _, importedComponent := range p.Imports() {
			imports = append(imports, NewImport(p.Component, importedComponent))
		}
	}

	g := NewGraph(imports)

	return g
}
