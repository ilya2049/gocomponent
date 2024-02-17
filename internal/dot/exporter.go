package dot

import (
	"strings"

	"github.com/ilya2049/gocomponent/internal/component"
)

type exporter struct {
	nameSpaceColors map[string]string
}

func newExporter(nameSpaceColors map[string]string) *exporter {
	return &exporter{
		nameSpaceColors: nameSpaceColors,
	}
}

func (*exporter) export(g *component.Graph) string {
	sb := strings.Builder{}

	sb.WriteString("digraph {\n")

	for _, component := range g.Components() {
		sb.WriteString(`"` + component.ID() + `"` + "\n")
	}

	for _, imp := range g.Imports() {
		sb.WriteString(`"` + imp.From().ID() + `" -> "` + imp.To().ID() + `"` + "\n")
	}

	sb.WriteString("}")

	return sb.String()
}
