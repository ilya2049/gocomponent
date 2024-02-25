package cliapp

import (
	"fmt"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/fs"
	"github.com/ilya2049/gocomponent/internal/project"
)

func PrintNamespaces() error {
	conf, err := config.Read()
	if err != nil {
		return err
	}

	prj := project.New()

	fsWalker := fs.NewWalk(conf.ProjectDirectory, prj)

	componentGraph, err := component.GenerateGraph(conf.ToComponentGraphConfig(), fsWalker)
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
