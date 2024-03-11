package main

import (
	"fmt"
	"os"

	"github.com/ilya2049/gocomponent/internal/cliapp"
	"github.com/ilya2049/gocomponent/internal/dot"
	"github.com/ilya2049/gocomponent/internal/fs"

	"github.com/urfave/cli/v2"
)

func main() {
	if err := newApp().Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func newApp() *cli.App {
	return &cli.App{
		Name:  "gocomponent",
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
				Name:   "svg",
				Usage:  "Print an svg file with a graph",
				Action: printDotSVG,
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

	server := cliapp.NewHTTPServer(serverPort, &fs.ComponentGraphReader{}, &dot.Exporter{})

	fmt.Println("Server started at " + server.Addr)

	return server.ListenAndServe()
}

func printDotGraph(*cli.Context) error {
	componentGraphReader := &fs.ComponentGraphReader{}

	componentGraph, err := componentGraphReader.ReadComponentGraph()
	if err != nil {
		return err
	}

	printer := cliapp.NewDotGraphPrinter(&dot.Exporter{}, os.Stdout)

	return printer.PrintDotGraph(componentGraph)
}

func printDotSVG(*cli.Context) error {
	componentGraphReader := &fs.ComponentGraphReader{}

	componentGraph, err := componentGraphReader.ReadComponentGraph()
	if err != nil {
		return err
	}

	printer := cliapp.NewDotSVGPrinter(&dot.Exporter{}, os.Stdout)

	return printer.PrintDotSVG(componentGraph)
}

func printNamespaces(*cli.Context) error {
	componentGraphReader := &fs.ComponentGraphReader{}

	componentGraph, err := componentGraphReader.ReadComponentGraph()
	if err != nil {
		return err
	}

	printer := cliapp.NewNamespacePrinter(os.Stdout)

	return printer.PrintNamespaces(componentGraph)
}
