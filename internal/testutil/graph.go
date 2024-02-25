package testutil

import (
	"strings"

	"github.com/ilya2049/gocomponent/internal/component"
)

func BuildGraphString(imports ...string) string {
	sb := strings.Builder{}

	for _, imp := range imports {
		sb.WriteString(imp + "\n")
	}

	return sb.String()
}

func GetComponentIDs(g *component.Graph) []string {
	var componentIDs []string

	for _, c := range g.Components() {
		componentIDs = append(componentIDs, c.ID())
	}

	return componentIDs
}

func AreComponentsEqual(component1, component2 *component.Component) bool {
	return component1.Namespace().String() == component2.Namespace().String()
}

type FsWalkerStub struct {
	componentGraph *component.Graph
}

func NewFsWalkerStub(componentGraph *component.Graph) *FsWalkerStub {
	return &FsWalkerStub{
		componentGraph: componentGraph,
	}
}

func (w *FsWalkerStub) CreateComponentGraph() (*component.Graph, error) {
	return w.componentGraph, nil
}
