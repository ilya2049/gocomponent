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

func (p *Project) Components() []*Component {
	var components []*Component

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

func (p *Project) ExcludeThirdPartyImports() {
	for _, aPackage := range p.packages {
		for _, packageImport := range aPackage.imports {
			if packageImport.isThirdParty {
				delete(aPackage.imports, packageImport.namespace)
			}
		}
	}
}

func (p *Project) IncludeOnlyNextPackageNamespaces(selectedPackageNamespaces []string) {
	filteredPackages := make(map[Namespace]*Package)

	for namespace, aPackage := range p.packages {
		for _, selectedNamespace := range selectedPackageNamespaces {
			if namespace.Contains(NewNamespace(selectedNamespace)) {
				filteredPackages[namespace] = aPackage
			}
		}
	}

	p.packages = filteredPackages
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
