package fswalk

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Walk struct {
	projectDir                    string
	rootNamespace                 string
	componentsAndTheirConnections map[string]map[string]struct{}
}

func New(projectDir string, rootNamespace string) *Walk {
	return &Walk{
		projectDir:                    projectDir,
		rootNamespace:                 rootNamespace,
		componentsAndTheirConnections: map[string]map[string]struct{}{},
	}
}

func (w *Walk) FindComponents() error {
	componentConnectionRegexp, err := regexp.Compile(
		fmt.Sprintf(componentConnectionRegexpTemplate, w.rootNamespace),
	)

	if err != nil {
		return fmt.Errorf("compile regexp to find component connections: %w", err)
	}

	startWalkHere := w.projectDir + "/" + w.rootNamespace + "/"

	return filepath.Walk(startWalkHere, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isGoSourceFile(path) {
			if component := getComponentUsingRegexp(path); component != "" {
				w.noteNewComponent(component)

				goFileContents, err := readFile(path)
				if err != nil {
					return err
				}

				for _, componentConnection := range getComponentConnectionsUsingRegexp(
					componentConnectionRegexp, goFileContents,
				) {
					w.noteNewComponentConnection(component, componentConnection)
				}
			}
		}

		return nil
	})
}

func (w *Walk) noteNewComponent(component string) {
	if _, ok := w.componentsAndTheirConnections[component]; !ok {
		w.componentsAndTheirConnections[component] = map[string]struct{}{}
	}
}

func (w *Walk) noteNewComponentConnection(component, componentConnection string) {
	if component == componentConnection {
		return
	}

	w.componentsAndTheirConnections[component][componentConnection] = struct{}{}
}

func (w *Walk) PrintDotGraph() {
	fmt.Println("digraph G {")

	for component, connections := range w.componentsAndTheirConnections {
		fmt.Println(component)

		if len(connections) > 0 {
			for connection := range connections {
				fmt.Println(component, "->", connection)
			}
		}
	}

	fmt.Println("}")
}
