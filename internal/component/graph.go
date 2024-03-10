package component

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type GraphConfig struct {
	ExtendComponentIDs          map[string]int
	IncludeThirdPartyComponents bool
	ThirdPartyComponentsColor   Color
	IncludeParentComponents     Namespaces
	IncludeChildComponents      Namespaces
	ExcludeParentComponents     Namespaces
	ExcludeChildComponents      Namespaces
	CustomComponents            Namespaces
	OnlyComponents              Namespaces
	ComponentColors             map[Namespace]Color
	EnableComponentSize         bool
}

func ApplyGraphConfig(conf *GraphConfig, componentGraph *Graph) (*Graph, error) {
	if !conf.IncludeThirdPartyComponents {
		componentGraph = componentGraph.RemoveThirdPartyComponents()
	}

	if len(conf.IncludeParentComponents) > 0 && len(conf.IncludeChildComponents) > 0 {
		componentGraph = componentGraph.IncludeParentAndChildComponents(
			conf.IncludeParentComponents, conf.IncludeChildComponents,
		)
	} else {
		if len(conf.IncludeParentComponents) > 0 {
			componentGraph = componentGraph.IncludeParentComponents(conf.IncludeParentComponents)
		}

		if len(conf.IncludeChildComponents) > 0 {
			componentGraph = componentGraph.IncludeChildComponents(conf.IncludeChildComponents)
		}
	}

	if len(conf.ExcludeParentComponents) > 0 {
		componentGraph = componentGraph.ExcludeParentComponents(
			conf.ExcludeParentComponents,
		)
	}

	if len(conf.ExcludeChildComponents) > 0 {
		componentGraph = componentGraph.ExcludeChildComponents(
			conf.ExcludeChildComponents,
		)
	}

	if len(conf.CustomComponents) > 0 {
		componentGraph = componentGraph.CreateCustomComponents(
			conf.CustomComponents,
		)
	}

	if len(conf.OnlyComponents) > 0 {
		componentGraph = componentGraph.IncludeOnlyComponents(conf.OnlyComponents)
	}

	if len(conf.ComponentColors) > 0 {
		componentGraph.Colorize(conf.ComponentColors)
	}

	if conf.ThirdPartyComponentsColor != "" {
		componentGraph.ColorizeThirdParty(conf.ThirdPartyComponentsColor)
	}

	componentGraph.MakeUniqueComponentIDs()

	if conf.EnableComponentSize {
		componentGraph.NormalizeComponentSizes()
	}

	if len(conf.ExtendComponentIDs) > 0 {
		if err := componentGraph.ExtendComponentIDs(conf.ExtendComponentIDs); err != nil {
			return nil, err
		}
	}

	return componentGraph, nil
}

type Graph struct {
	components map[Namespace]*Component

	imports Imports
}

func NewGraph(imports Imports) *Graph {
	componentsMap := make(map[Namespace]*Component)
	importsMap := make(map[string]*Import)
	newImports := make(Imports, 0)

	for _, imp := range imports {
		if _, ok := importsMap[imp.String()]; ok {
			continue
		}

		importsMap[imp.String()] = imp
		newImports = append(newImports, imp)

		componentsMap[imp.from.namespace] = imp.from
		componentsMap[imp.to.namespace] = imp.to
	}

	return &Graph{
		components: componentsMap,

		imports: newImports,
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

func (g *Graph) NormalizeComponentSizes() {
	const minNormalizeSize = 0.1

	components := g.Components()

	if len(components) == 0 {
		return
	}

	componentWithMaxSize := slices.MaxFunc(components, func(a, b *Component) int {
		return a.SizeBytes() - b.SizeBytes()
	})

	componentWithMinSize := slices.MinFunc(components, func(a, b *Component) int {
		return a.SizeBytes() - b.SizeBytes()
	})

	maxMinDifference := componentWithMaxSize.SizeBytes() - componentWithMinSize.SizeBytes()

	for _, component := range components {
		normalizedSize := float64((component.SizeBytes() - componentWithMinSize.SizeBytes())) / float64(maxMinDifference)

		if normalizedSize < minNormalizeSize {
			normalizedSize = minNormalizeSize
		}

		component.NormalizeSize(normalizedSize)
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
				if sections <= 0 {
					component.UseNamespaceAsID()

					continue
				}

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

func (g *Graph) IncludeParentAndChildComponents(
	parentComponentNamespaces Namespaces,
	childComponentNamespaces Namespaces,
) *Graph {
	parentComponentImports := g.includeParentComponentImports(parentComponentNamespaces)
	childComponentImports := g.includeChildComponentImports(childComponentNamespaces)

	return NewGraph(append(parentComponentImports, childComponentImports...))
}

func (g *Graph) IncludeParentComponents(namespaces Namespaces) *Graph {
	return NewGraph(g.includeParentComponentImports(namespaces))
}

func (g *Graph) includeParentComponentImports(namespaces Namespaces) Imports {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		includeImport := false

		for _, namespace := range namespaces {
			if imp.from.namespace.Contains(namespace) {
				includeImport = true

				break
			}
		}

		if includeImport {
			newImports = append(newImports, imp)
		}
	}

	return newImports
}

func (g *Graph) IncludeChildComponents(namespaces Namespaces) *Graph {
	return NewGraph(g.includeChildComponentImports(namespaces))
}

func (g *Graph) includeChildComponentImports(namespaces Namespaces) Imports {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		includeImport := false

		for _, namespace := range namespaces {
			if imp.to.namespace.Contains(namespace) {
				includeImport = true

				break
			}
		}

		if includeImport {
			newImports = append(newImports, imp)
		}
	}

	return newImports
}

func (g *Graph) ExcludeParentComponents(namespaces Namespaces) *Graph {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		includeImport := true

		for _, namespace := range namespaces {
			if imp.from.namespace.Contains(namespace) {
				includeImport = false

				break
			}
		}

		if includeImport {
			newImports = append(newImports, imp)
		}
	}

	return NewGraph(newImports)
}

func (g *Graph) ExcludeChildComponents(namespaces Namespaces) *Graph {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		includeImport := true

		for _, namespace := range namespaces {
			if imp.to.namespace.Contains(namespace) {
				includeImport = false

				break
			}
		}

		if includeImport {
			newImports = append(newImports, imp)
		}
	}

	return NewGraph(newImports)
}

func (g *Graph) IncludeOnlyComponents(namespaces Namespaces) *Graph {
	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		fromComponentInNamespace := false
		toComponentInNamespace := false

		for _, namespace := range namespaces {
			if imp.From().Namespace().Contains(namespace) {
				fromComponentInNamespace = true
			}

			if imp.To().Namespace().Contains(namespace) {
				toComponentInNamespace = true
			}
		}

		if fromComponentInNamespace && toComponentInNamespace {
			newImports = append(newImports, imp)
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
	componentsToMergeInCustomComponent := make(map[Namespace]*Component)

	newImports := make(Imports, 0)

	for _, imp := range g.Imports() {
		if imp.from.namespace.Contains(namespace) && imp.to.namespace.Contains(namespace) {
			continue
		}

		if imp.from.namespace.Contains(namespace) {
			componentsToMergeInCustomComponent[imp.from.namespace] = imp.from

			childrenOfCustomComponent[imp.to.namespace] = imp.to

			continue
		}

		if imp.to.namespace.Contains(namespace) {
			componentsToMergeInCustomComponent[imp.to.namespace] = imp.to

			parentsOfCustomComponent[imp.from.namespace] = imp.from

			continue
		}

		newImports = append(newImports, imp)
	}

	customComponent := New(namespace)

	for _, componentToMergeInCustomComponent := range componentsToMergeInCustomComponent {
		customComponent.AddBytesInSize(componentToMergeInCustomComponent.SizeBytes())
	}

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

	graphImports := g.Imports()

	for i, imp := range graphImports {
		sb.WriteString(imp.String())

		if i < len(graphImports)-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}
