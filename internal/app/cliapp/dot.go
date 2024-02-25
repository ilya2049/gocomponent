package cliapp

import (
	"fmt"

	"github.com/ilya2049/gocomponent/internal/domain/component"
)

type dotExporter interface {
	Export(*component.Graph) string
}

type dotSVGExporter interface {
	ExportSVG(*component.Graph) ([]byte, error)
}

func PrintDotGraph(
	conf *component.GraphConfig,
	initialComponentGraph *component.Graph,
	dotExporter dotExporter,
) error {
	componentGraph, err := component.ApplyGraphConfig(conf, initialComponentGraph)
	if err != nil {
		return err
	}

	fmt.Println(dotExporter.Export(componentGraph))

	return nil
}
