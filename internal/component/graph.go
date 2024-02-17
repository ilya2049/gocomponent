package component

type Graph struct {
	components Components
	imports    Imports
}

func NewGraph(components Components, imports Imports) *Graph {
	return &Graph{
		components: components,
		imports:    imports,
	}
}

func (g *Graph) Components() Components {
	return g.components
}

func (g *Graph) Imports() Imports {
	return g.imports
}

func (g *Graph) RemoveThirdPartyComponents() *Graph {
	newComponents := make(Components, 0)

	for _, component := range g.Components() {
		if !component.isThirdParty {
			newComponents = append(newComponents, component)
		}
	}

	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		if !imp.to.isThirdParty {
			newImports = append(newImports, imp)
		}
	}

	return NewGraph(newComponents, newImports)
}
