package cliapp

import (
	"fmt"
	"io"

	"github.com/ilya2049/gocomponent/internal/component"
)

type dotExporter interface {
	Export(*component.Graph) string
}

type dotSVGExporter interface {
	ExportSVG(*component.Graph) ([]byte, error)
}

type DotGraphPrinter struct {
	dotExporter dotExporter
	destination io.Writer
}

func NewDotGraphPrinter(dotExporter dotExporter, destination io.Writer) *DotGraphPrinter {
	return &DotGraphPrinter{
		dotExporter: dotExporter,
		destination: destination,
	}
}

func (p *DotGraphPrinter) PrintDotGraph(componentGraph *component.Graph) error {
	fmt.Fprintln(p.destination, p.dotExporter.Export(componentGraph))

	return nil
}
