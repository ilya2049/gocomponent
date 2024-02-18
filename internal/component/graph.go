package component

type Graph struct {
	components map[Namespace]*Component

	imports Imports
}

func NewGraph(imports Imports) *Graph {
	componentsMap := make(map[Namespace]*Component)

	for _, imp := range imports {
		componentsMap[imp.from.namespace] = imp.from
		componentsMap[imp.to.namespace] = imp.to
	}

	return &Graph{
		components: componentsMap,

		imports: imports,
	}
}

func (g *Graph) Components() Components {
	components := make(Components, 0, len(g.components))

	for _, component := range g.components {
		components = append(components, component)
	}

	return components
}

func (g *Graph) Imports() Imports {
	return g.imports
}

func (g *Graph) Colorize(namespaceColorMap map[Namespace]Color) {
	for _, component := range g.components {
		for namespace, color := range namespaceColorMap {
			if component.namespace.Contains(namespace) {
				component.Colorize(color)

				continue
			}
		}
	}
}

func (g *Graph) ColorizeThirdParty(color Color) {
	for _, component := range g.components {
		if component.isThirdParty {
			component.Colorize(color)
		}
	}
}

func (g *Graph) RemoveThirdPartyComponents() *Graph {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		if !imp.to.isThirdParty {
			newImports = append(newImports, imp)
		}
	}

	return NewGraph(newImports)
}

func (g *Graph) IncludeParentComponents(namespaces Namespaces) *Graph {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		for _, namespace := range namespaces {
			if imp.from.namespace.Contains(namespace) {
				newImports = append(newImports, imp)

				continue
			}
		}
	}

	return NewGraph(newImports)
}

func (g *Graph) IncludeChildComponents(namespaces Namespaces) *Graph {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		for _, namespace := range namespaces {
			if imp.to.namespace.Contains(namespace) {
				newImports = append(newImports, imp)

				continue
			}
		}
	}

	return NewGraph(newImports)
}

func (g *Graph) ExcludeParentComponents(namespaces Namespaces) *Graph {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		for _, namespace := range namespaces {
			if !imp.from.namespace.Contains(namespace) {
				newImports = append(newImports, imp)

				continue
			}
		}
	}

	return NewGraph(newImports)
}

func (g *Graph) ExcludeChildComponents(namespaces Namespaces) *Graph {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		for _, namespace := range namespaces {
			if !imp.to.namespace.Contains(namespace) {
				newImports = append(newImports, imp)

				continue
			}
		}
	}

	return NewGraph(newImports)
}
