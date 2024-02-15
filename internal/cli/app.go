package cli

import (
	"fmt"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/dot"
	"github.com/ilya2049/gocomponent/internal/fs"
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
	graph := dot.GenerateGraph()

	fmt.Println(graph)

	return nil
}

func printNamespaces(cCtx *cli.Context) error {
	conf, err := config.Read()
	if err != nil {
		return err
	}

	project := component.NewProject()

	walk := fs.NewWalk(conf.ProjectDirectory, project)

	if err := walk.FindComponentsAndImports(); err != nil {
		return err
	}

	for _, pkg := range project.Packages() {
		fmt.Println(pkg.Namespace())
	}

	return nil
}
