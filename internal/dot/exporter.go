package dot

import (
	"strings"

	"github.com/ilya2049/gocomponent/internal/component"
)

type Exporter struct {
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Export(packages []*component.Package) string {
	sb := strings.Builder{}

	sb.WriteString(e.startGraph())

	for _, p := range packages {
		sb.WriteString(e.addPackageInGraph(p))

		for _, importedComponent := range p.Imports() {
			sb.WriteString(e.addImportInGraph(p, importedComponent))
		}
	}

	sb.WriteString(e.completeGraph())

	return sb.String()
}

func (*Exporter) startGraph() string {
	return "digraph {\n"
}

func (*Exporter) completeGraph() string {
	return "}"
}

func (*Exporter) addPackageInGraph(aPackage *component.Package) string {
	return `"` + aPackage.ID() + `"` + "\n"
}

func (*Exporter) addImportInGraph(aPackage *component.Package, importedComponent *component.Component) string {
	return `"` + aPackage.ID() + `" -> "` + importedComponent.ID() + `"` + "\n"
}
