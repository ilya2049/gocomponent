package cliapp

import (
	"fmt"
	"io"
	"strings"

	"github.com/ilya2049/gocomponent/internal/component"
)

type NamespacePrinter struct {
	destination io.Writer
}

func NewNamespacePrinter(destination io.Writer) *NamespacePrinter {
	return &NamespacePrinter{
		destination: destination,
	}
}

func (p *NamespacePrinter) PrintNamespaces(componentGraph *component.Graph) error {
	sb := strings.Builder{}

	components := componentGraph.Components().OrderByStability()

	for i, component := range components {
		sb.WriteString(fmt.Sprintf("%.2f %s [%s]", component.Stability(), component.Namespace(), component.ID()))

		if i < len(components)-1 {
			sb.WriteRune('\n')
		}
	}

	fmt.Fprintln(p.destination, sb.String())

	return nil
}
