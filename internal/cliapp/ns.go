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

	initialComponentGraph, err := fsWalker.CreateComponentGraph()
	if err != nil {
		return err
	}

	componentGraph, err := component.ApplyGraphConfig(conf.ToComponentGraphConfig(), initialComponentGraph)
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
