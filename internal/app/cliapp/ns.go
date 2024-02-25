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

func (p *NamespacePrinter) PrintNamespaces(conf *component.GraphConfig, initialComponentGraph *component.Graph) error {
	componentGraph, err := component.ApplyGraphConfig(conf, initialComponentGraph)
	if err != nil {
		return err
	}

	sb := strings.Builder{}

	components := componentGraph.Components()

	for i, component := range components {
		sb.WriteString(fmt.Sprintf("%s [%s]", component.Namespace(), component.ID()))

		if i < len(components)-1 {
			sb.WriteRune('\n')
		}
	}

	fmt.Fprintln(p.destination, sb.String())

	return nil
}
