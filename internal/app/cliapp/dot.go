package cliapp

import (
	"fmt"
	"io"

	"github.com/ilya2049/gocomponent/internal/domain/component"
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

func (p *DotGraphPrinter) PrintDotGraph(conf *component.GraphConfig, initialComponentGraph *component.Graph) error {
	componentGraph, err := component.ApplyGraphConfig(conf, initialComponentGraph)
	if err != nil {
		return err
	}

	fmt.Fprintln(p.destination, p.dotExporter.Export(componentGraph))

	return nil
}
