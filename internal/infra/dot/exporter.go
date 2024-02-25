package dot

import (
	"bytes"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/ilya2049/gocomponent/internal/domain/component"
)

type Exporter struct {
}

func (e *Exporter) Export(g *component.Graph) string {
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

func (e *Exporter) ExportSVG(g *component.Graph) ([]byte, error) {
	dotGraph := e.Export(g)

	parsedDotGraph, err := graphviz.ParseBytes([]byte(dotGraph))
	if err != nil {
		return nil, err
	}

	var svgGraph bytes.Buffer

	graph := graphviz.New()
	if err := graph.Render(parsedDotGraph, graphviz.SVG, &svgGraph); err != nil {
		return nil, err
	}

	return mapComponentAndNamespaceInSVG(g.Components(), svgGraph.Bytes()), nil
}

func mapComponentAndNamespaceInSVG(components component.Components, svgBytes []byte) []byte {
	svgString := string(svgBytes)

	for _, aComponent := range components {
		svgString = strings.Replace(svgString,
			"<title>"+aComponent.ID()+"</title>",
			"<title>"+aComponent.Namespace().String()+"</title>",
			1,
		)
	}

	return []byte(svgString)
}
