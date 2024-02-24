package component

import (
	"fmt"
	"regexp"
	"strings"
)

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

func (g *Graph) MakeUniqueComponentIDs() {
	components := g.Components()

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

func (g *Graph) ExtendComponentIDs(idRegexpPatternAndSections map[string]int) error {
	for idRegexpPattern, sections := range idRegexpPatternAndSections {
		r, err := regexp.Compile(idRegexpPattern)
		if err != nil {
			return fmt.Errorf("extend component ids: %w", err)
		}

		for _, component := range g.components {
			if r.MatchString(component.Namespace().String()) {
				for i := 0; i < sections; i++ {
					component.ExtendID()
				}
			}
		}
	}

	return nil
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
		if component.IsThirdParty() {
			component.Colorize(color)
		}
	}
}

func (g *Graph) RemoveThirdPartyComponents() *Graph {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		if !imp.to.IsThirdParty() {
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

func (g *Graph) CreateCustomComponents(namespaces Namespaces) *Graph {
	newGraph := g

	for _, namespace := range namespaces {
		newGraph = newGraph.createCustomComponent(namespace)
	}

	return newGraph
}

func (g *Graph) createCustomComponent(namespace Namespace) *Graph {
	childrenOfCustomComponent := make(map[Namespace]*Component)
	parentsOfCustomComponent := make(map[Namespace]*Component)

	newImports := make(Imports, 0)

	for _, imp := range g.imports {
		if imp.from.namespace.Contains(namespace) && imp.to.namespace.Contains(namespace) {
			continue
		}

		if imp.from.namespace.Contains(namespace) {
			childrenOfCustomComponent[imp.to.namespace] = imp.to

			continue
		}

		if imp.to.namespace.Contains(namespace) {
			parentsOfCustomComponent[imp.from.namespace] = imp.from

			continue
		}

		newImports = append(newImports, imp)
	}

	customComponent := New(namespace)

	for _, childOfCustomComponent := range childrenOfCustomComponent {
		newImports = append(newImports, NewImport(customComponent, childOfCustomComponent))
	}

	for _, parentOfCustomComponent := range parentsOfCustomComponent {
		newImports = append(newImports, NewImport(parentOfCustomComponent, customComponent))
	}

	return NewGraph(newImports)
}

func (g *Graph) String() string {
	sb := strings.Builder{}

	for _, imp := range g.Imports() {
		sb.WriteString(imp.From().Namespace().String() + " -> " + imp.To().Namespace().String() + "\n")
	}

	return sb.String()
}
