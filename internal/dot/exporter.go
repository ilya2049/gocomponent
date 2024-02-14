package dot

import (
	"strings"

	"github.com/ilya2049/gocomponent/internal/component"
)

type exporter struct {
}

func newExporter() *exporter {
	return &exporter{}
}

func (e *exporter) export(packages []*component.Package) string {
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

func (*exporter) startGraph() string {
	return "digraph {\n"
}

func (*exporter) completeGraph() string {
	return "}"
}

func (*exporter) addPackageInGraph(aPackage *component.Package) string {
	return `"` + aPackage.ID() + `"` + "\n"
}

func (*exporter) addImportInGraph(aPackage *component.Package, importedComponent *component.Component) string {
	return `"` + aPackage.ID() + `" -> "` + importedComponent.ID() + `"` + "\n"
}
