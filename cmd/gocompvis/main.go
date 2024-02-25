package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilya2049/gocomponent/internal/cliapp"
	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
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

	server := cliapp.NewHTTPServer(serverPort, readComponentGraph)

	fmt.Println("Server started at " + server.Addr)

	return server.ListenAndServe()
}

func printDotGraph(*cli.Context) error {
	conf, initialComponentGraph, err := readComponentGraph()
	if err != nil {
		return err
	}

	return cliapp.PrintDotGraph(conf, initialComponentGraph)
}

func printNamespaces(*cli.Context) error {
	conf, initialComponentGraph, err := readComponentGraph()
	if err != nil {
		return err
	}

	return cliapp.PrintNamespaces(conf, initialComponentGraph)
}

func readComponentGraph() (*component.GraphConfig, *component.Graph, error) {
	conf, err := config.Read()
	if err != nil {
		return nil, nil, err
	}

	fsWalker := fs.NewWalk(conf.ProjectDirectory, &fileReader{}, &filePathWalker{})

	initialComponentGraph, err := fsWalker.ReadComponentGraph()
	if err != nil {
		return nil, nil, err
	}

	return conf.ToComponentGraphConfig(), initialComponentGraph, nil
}

type fileReader struct {
}

func (r *fileReader) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

type filePathWalker struct {
}

func (w *filePathWalker) Walk(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(root, fn)
}
