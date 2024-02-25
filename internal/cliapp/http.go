package cliapp

import (
	"bytes"
	"net/http"

	"github.com/goccy/go-graphviz"
	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/dot"
	"github.com/ilya2049/gocomponent/internal/fs"
	"github.com/ilya2049/gocomponent/internal/project"
)

const defaultHTTPServerPort = "8080"

func NewHTTPServer(port string) *http.Server {
	if port == "" {
		port = defaultHTTPServerPort
	}

	server := http.Server{
		Addr: ":" + port,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHTTPRequest)

	server.Handler = mux

	return &server
}

func handleHTTPRequest(w http.ResponseWriter, _ *http.Request) {
	conf, err := config.Read()
	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	prj := project.New()

	fsWalker := fs.NewWalk(conf.ProjectDirectory, prj)

	componentGraph, err := component.GenerateGraph(conf.ToComponentGraphConfig(), fsWalker)
	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	dotGraph := dot.Export(componentGraph)

	parsedDotGraph, err := graphviz.ParseBytes([]byte(dotGraph))
	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	var svgGraph bytes.Buffer

	graph := graphviz.New()
	if err := graph.Render(parsedDotGraph, graphviz.SVG, &svgGraph); err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	w.Write(dot.MapComponentAndNamespaceInSVG(componentGraph.Components(), svgGraph.Bytes()))
}
