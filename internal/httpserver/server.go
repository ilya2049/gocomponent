package httpserver

import (
	"bytes"
	"net/http"

	"github.com/goccy/go-graphviz"
	"github.com/ilya2049/gocomponent/internal/dot"
	"github.com/ilya2049/gocomponent/internal/generator"
)

func New(address string) *http.Server {
	server := http.Server{
		Addr: address,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		componentGraph, err := generator.GenerateGraph()
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

		w.Write(dot.MapComponentAndNamespace(componentGraph.Components(), svgGraph.Bytes()))
	})

	server.Handler = mux

	return &server
}
