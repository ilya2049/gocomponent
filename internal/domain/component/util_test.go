package component_test

import (
	"github.com/ilya2049/gocomponent/internal/domain/component"
)

func getComponentIDs(g *component.Graph) []string {
	var componentIDs []string

	for _, c := range g.Components() {
		componentIDs = append(componentIDs, c.ID())
	}

	return componentIDs
}

func areComponentsEqual(component1, component2 *component.Component) bool {
	return component1.Namespace().String() == component2.Namespace().String()
}
