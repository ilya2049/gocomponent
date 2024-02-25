package httpserver

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

func New(address string) *http.Server {
	server := http.Server{
		Addr: address,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	server.Handler = mux

	return &server
}
