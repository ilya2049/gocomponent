package fs

import "github.com/ilya2049/gocomponent/internal/domain/component"

type project struct {
	packages map[component.Namespace]*projectPackage

	components map[component.Namespace]*component.Component
}

func newProject() *project {
	return &project{
		packages:   make(map[component.Namespace]*projectPackage),
		components: make(map[component.Namespace]*component.Component),
	}
}

func (p *project) getOrAddComponent(namespace component.Namespace) *component.Component {
	existingComponent, ok := p.components[namespace]
	if ok {
		return existingComponent
	}

	newComponent := component.New(namespace)

	p.components[namespace] = newComponent

	return newComponent
}

func (p *project) findPackage(namespace component.Namespace) (*projectPackage, bool) {
	pkg, ok := p.packages[namespace]

	return pkg, ok
}

func (p *project) addPackage(namespace component.Namespace, pkg *projectPackage) {
	p.packages[namespace] = pkg
}

func (p *project) createComponentGraph() *component.Graph {
	imports := make(component.Imports, 0)

	for _, p := range p.packages {
		for _, importedComponent := range p.getImports() {
			imports = append(imports, component.NewImport(p.Component, importedComponent))
		}
	}

	g := component.NewGraph(imports)

	return g
}
