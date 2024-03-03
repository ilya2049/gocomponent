package cliapp

import (
	"net/http"
	"strings"

	"github.com/ilya2049/gocomponent/internal/domain/component"
)

type componentGraphReader interface {
	ReadComponentGraph(configFileName string) (*component.Graph, error)
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

	handler := newHTTPRequestHandler(componentGraphReader, dotSVGExporter)

	mux.HandleFunc("/", handler.handle)

	server.Handler = mux

	return &server
}

type httpRequestHandler struct {
	componentGraphReader componentGraphReader
	dotSVGExporter       dotSVGExporter
}

func newHTTPRequestHandler(
	componentGraphReader componentGraphReader,
	dotSVGExporter dotSVGExporter,
) *httpRequestHandler {
	return &httpRequestHandler{
		componentGraphReader: componentGraphReader,
		dotSVGExporter:       dotSVGExporter,
	}
}

func (h *httpRequestHandler) handle(w http.ResponseWriter, r *http.Request) {
	componentGraph, err := h.componentGraphReader.ReadComponentGraph(h.getConfigFileName(r))
	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	dotSVGGraph, err := h.dotSVGExporter.ExportSVG(componentGraph)
	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	w.Write(dotSVGGraph)
}

func (h *httpRequestHandler) getConfigFileName(r *http.Request) string {
	requestURISections := strings.Split(r.RequestURI, "/")

	var configFileName string

	if secondRequestURISection := requestURISections[1]; secondRequestURISection != "" {
		configFileName = secondRequestURISection + ".toml"
	}

	return configFileName
}
