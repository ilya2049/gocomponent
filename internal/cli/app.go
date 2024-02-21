package cli

import (
	"fmt"

	"github.com/ilya2049/gocomponent/internal/dot"
	"github.com/ilya2049/gocomponent/internal/generator"
	"github.com/ilya2049/gocomponent/internal/httpserver"

	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
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
	if serverPort == "" {
		serverPort = "8080"
	}

	server := httpserver.New(":" + serverPort)

	fmt.Println("Server started at " + serverPort)

	return server.ListenAndServe()
}

func printDotGraph(cCtx *cli.Context) error {
	componentGraph, err := generator.GenerateGraph()
	if err != nil {
		return err
	}

	fmt.Println(dot.Export(componentGraph))

	return nil
}

func printNamespaces(cCtx *cli.Context) error {
	componentGraph, err := generator.GenerateGraph()
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
