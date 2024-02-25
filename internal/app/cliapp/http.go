package cliapp

import (
	"bytes"
	"net/http"

	"github.com/goccy/go-graphviz"
	"github.com/ilya2049/gocomponent/internal/app/dot"
	"github.com/ilya2049/gocomponent/internal/domain/component"
)

type readComponentGraphFunc func() (*component.GraphConfig, *component.Graph, error)

const defaultHTTPServerPort = "8080"

func NewHTTPServer(port string, readComponentGraph readComponentGraphFunc) *http.Server {
	if port == "" {
		port = defaultHTTPServerPort
	}

	server := http.Server{
		Addr: ":" + port,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHTTPRequest(readComponentGraph))

	server.Handler = mux

	return &server
}

func handleHTTPRequest(readComponentGraph readComponentGraphFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		conf, initialComponentGraph, err := readComponentGraph()
		if err != nil {
			w.Write([]byte(err.Error()))

			return
		}

		componentGraph, err := component.ApplyGraphConfig(conf, initialComponentGraph)
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
}
