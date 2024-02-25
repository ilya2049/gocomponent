package cliapp

import (
	"net/http"

	"github.com/ilya2049/gocomponent/internal/domain/component"
)

type componentGraphReader interface {
	ReadComponentGraph() (*component.GraphConfig, *component.Graph, error)
}

const defaultHTTPServerPort = "8080"

func NewHTTPServer(
	port string,
	componentGraphReader componentGraphReader,
	dotSVGExporter dotSVGExporter,
) *http.Server {
	if port == "" {
		port = defaultHTTPServerPort
	}

	server := http.Server{
		Addr: ":" + port,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHTTPRequest(componentGraphReader, dotSVGExporter))

	server.Handler = mux

	return &server
}

func handleHTTPRequest(
	componentGraphReader componentGraphReader,
	dotSVGExporter dotSVGExporter,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		conf, initialComponentGraph, err := componentGraphReader.ReadComponentGraph()
		if err != nil {
			w.Write([]byte(err.Error()))

			return
		}

		componentGraph, err := component.ApplyGraphConfig(conf, initialComponentGraph)
		if err != nil {
			w.Write([]byte(err.Error()))

			return
		}

		dotSVGGraph, err := dotSVGExporter.ExportSVG(componentGraph)
		if err != nil {
			w.Write([]byte(err.Error()))

			return
		}

		w.Write(dotSVGGraph)
	}
}
