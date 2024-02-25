package main

import (
	"fmt"
	"os"

	"github.com/ilya2049/gocomponent/internal/cliapp"
	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/dot"
	"github.com/ilya2049/gocomponent/internal/fs"
	"github.com/ilya2049/gocomponent/internal/project"

	"github.com/urfave/cli/v2"
)

func main() {
	if err := newApp().Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func newApp() *cli.App {
	return &cli.App{
		Name:  "gocompvis",
		Usage: "Visualize your project components and their connections",
		Commands: []*cli.Command{
			{
				Name:   "http",
				Usage:  "Run an http server",
				Action: runHTTPServer,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Http server port",
					},
				},
			},
			{
				Name:   "dot",
				Usage:  "Print a dot graph",
				Action: printDotGraph,
			},
			{
				Name:   "ns",
				Usage:  "Print namespaces in the project",
				Action: printNamespaces,
			},
		},
	}
}

func runHTTPServer(cCtx *cli.Context) error {
	serverPort := cCtx.String("port")

	server := cliapp.NewHTTPServer(serverPort)

	fmt.Println("Server started at " + server.Addr)

	return server.ListenAndServe()
}

func printDotGraph(cCtx *cli.Context) error {
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

	fmt.Println(dot.Export(componentGraph))

	return nil
}

func printNamespaces(cCtx *cli.Context) error {
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
