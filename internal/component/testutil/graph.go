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
