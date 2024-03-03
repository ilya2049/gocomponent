package cliapp

import (
	"fmt"
	"io"
	"strings"

	"github.com/ilya2049/gocomponent/internal/domain/component"
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

	components := componentGraph.Components()

	for i, component := range components {
		sb.WriteString(fmt.Sprintf("%s [%s] %d B",
			component.Namespace(), component.ID(), component.SizeBytes(),
		))

		if i < len(components)-1 {
			sb.WriteRune('\n')
		}
	}

	fmt.Fprintln(p.destination, sb.String())

	return nil
}
