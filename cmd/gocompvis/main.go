package main

import (
	"fmt"
	"os"

	"github.com/ilya2049/gocomponent/internal/app/cliapp"
	"github.com/ilya2049/gocomponent/internal/infra/fs"

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

	server := cliapp.NewHTTPServer(serverPort, fs.ReadComponentGraph)

	fmt.Println("Server started at " + server.Addr)

	return server.ListenAndServe()
}

func printDotGraph(*cli.Context) error {
	conf, initialComponentGraph, err := fs.ReadComponentGraph()
	if err != nil {
		return err
	}

	return cliapp.PrintDotGraph(conf, initialComponentGraph)
}

func printNamespaces(*cli.Context) error {
	conf, initialComponentGraph, err := fs.ReadComponentGraph()
	if err != nil {
		return err
	}

	return cliapp.PrintNamespaces(conf, initialComponentGraph)
}
