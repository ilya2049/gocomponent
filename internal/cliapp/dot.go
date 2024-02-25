package cliapp

import (
	"fmt"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/dot"
)

func PrintDotGraph(conf *component.GraphConfig, initialComponentGraph *component.Graph) error {
	componentGraph, err := component.ApplyGraphConfig(conf, initialComponentGraph)
	if err != nil {
		return err
	}

	fmt.Println(dot.Export(componentGraph))

	return nil
}
