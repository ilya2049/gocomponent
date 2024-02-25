package cliapp

import (
	"fmt"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/dot"
	"github.com/ilya2049/gocomponent/internal/fs"
	"github.com/ilya2049/gocomponent/internal/project"
)

func PrintDotGraph() error {
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

	fmt.Println(dot.Export(componentGraph))

	return nil
}
