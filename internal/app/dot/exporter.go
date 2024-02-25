package dot

import (
	"strings"

	"github.com/ilya2049/gocomponent/internal/domain/component"
)

func Export(g *component.Graph) string {
	sb := strings.Builder{}

	sb.WriteString("digraph {\n")

	for _, component := range g.Components() {
		componentString := `"` + component.ID() + `"`

		if component.HasColor() {
			componentString += " [shape=component, style=filled, fillcolor=" +
				component.Color().String() + "]\n"
		} else {
			componentString += " [shape=component]\n"
		}

		sb.WriteString(componentString)
	}

	for _, imp := range g.Imports() {
		sb.WriteString(`"` + imp.From().ID() + `" -> "` + imp.To().ID() + `"` + "\n")
	}

	sb.WriteString("}")

	return sb.String()
}
