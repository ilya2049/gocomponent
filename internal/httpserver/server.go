package httpserver

import (
	"bytes"
	"net/http"

	"github.com/goccy/go-graphviz"
	"github.com/ilya2049/gocomponent/internal/dot"
)

func New(address string) *http.Server {
	server := http.Server{
		Addr: address,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		parsedDotGraph, err := graphviz.ParseBytes([]byte(dot.GenerateGraph()))
		if err != nil {
			w.Write([]byte(err.Error()))
		}

		var svgGraph bytes.Buffer

		graph := graphviz.New()
		if err := graph.Render(parsedDotGraph, graphviz.SVG, &svgGraph); err != nil {
			w.Write([]byte(err.Error()))
		}

		w.Write(svgGraph.Bytes())
	})

	server.Handler = mux

	return &server
}
