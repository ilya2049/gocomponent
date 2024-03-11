package cliapp

import (
	"fmt"
	"io"

	"github.com/ilya2049/gocomponent/internal/component"
)

type dotSVGExporter interface {
	ExportSVG(*component.Graph) ([]byte, error)
}

type DotSVGPrinter struct {
	dotSVGExporter dotSVGExporter
	destination    io.Writer
}

func NewDotSVGPrinter(dotSVGExporter dotSVGExporter, destination io.Writer) *DotSVGPrinter {
	return &DotSVGPrinter{
		dotSVGExporter: dotSVGExporter,
		destination:    destination,
	}
}

func (p *DotSVGPrinter) PrintDotSVG(componentGraph *component.Graph) error {
	svgBytes, err := p.dotSVGExporter.ExportSVG(componentGraph)
	if err != nil {
		return fmt.Errorf("export svg: %w", err)
	}

	fmt.Fprintln(p.destination, string(svgBytes))

	return nil
}
