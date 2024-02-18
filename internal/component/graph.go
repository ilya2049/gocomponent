package component

type Graph struct {
	componentsMap map[Namespace]*Component

	imports Imports
}

func NewGraph(imports Imports) *Graph {
	componentsMap := make(map[Namespace]*Component)

	for _, imp := range imports {
		componentsMap[imp.from.namespace] = imp.from
		componentsMap[imp.to.namespace] = imp.to
	}

	return &Graph{
		componentsMap: componentsMap,

		imports: imports,
	}
}

func (g *Graph) Components() Components {
	components := make(Components, 0, len(g.componentsMap))

	for _, component := range g.componentsMap {
		components = append(components, component)
	}

	return components
}

func (g *Graph) Imports() Imports {
	return g.imports
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
