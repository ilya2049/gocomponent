package project

import "github.com/ilya2049/gocomponent/internal/component"

type Project struct {
	packages map[component.Namespace]*Package

	components map[component.Namespace]*component.Component
}

func New() *Project {
	return &Project{
		packages:   make(map[component.Namespace]*Package),
		components: make(map[component.Namespace]*component.Component),
	}
}

func (p *Project) GetOrAddComponent(namespace component.Namespace) *component.Component {
	existingComponent, ok := p.components[namespace]
	if ok {
		return existingComponent
	}

	newComponent := component.New(namespace)

	p.components[namespace] = newComponent

	return newComponent
}

func (p *Project) Components() component.Components {
	var components component.Components

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

func (p *Project) FindPackage(namespace component.Namespace) (*Package, bool) {
	pkg, ok := p.packages[namespace]

	return pkg, ok
}

func (p *Project) AddPackage(namespace component.Namespace, pkg *Package) {
	p.packages[namespace] = pkg
}

func (p *Project) Packages() []*Package {
	packages := make([]*Package, 0, len(p.packages))

	for _, p := range p.packages {
		packages = append(packages, p)
	}

	return packages
}

func (p *Project) CreateComponentGraph() *component.Graph {
	imports := make(component.Imports, 0)

	for _, p := range p.packages {
		for _, importedComponent := range p.Imports() {
			imports = append(imports, component.NewImport(p.Component, importedComponent))
		}
	}

	g := component.NewGraph(imports)

	return g
}
