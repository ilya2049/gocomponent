package cliapp

import (
	"fmt"

	"github.com/ilya2049/gocomponent/internal/app/dot"
	"github.com/ilya2049/gocomponent/internal/domain/component"
)

func PrintDotGraph(conf *component.GraphConfig, initialComponentGraph *component.Graph) error {
	componentGraph, err := component.ApplyGraphConfig(conf, initialComponentGraph)
	if err != nil {
		return err
	}

	fmt.Println(dot.Export(componentGraph))

	return nil
}
