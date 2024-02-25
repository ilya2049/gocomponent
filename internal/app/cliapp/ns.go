package cliapp

import (
	"fmt"

	"github.com/ilya2049/gocomponent/internal/domain/component"
)

func PrintNamespaces(conf *component.GraphConfig, initialComponentGraph *component.Graph) error {
	componentGraph, err := component.ApplyGraphConfig(conf, initialComponentGraph)
	if err != nil {
		return err
	}

	for _, component := range componentGraph.Components() {
		var thirdParty string
		if component.IsThirdParty() {
			thirdParty = "*"
		}

		fmt.Println(component.Namespace(), "["+thirdParty+component.ID()+"]")
	}

	return nil
}
