package fswalk

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	// '/internal/app/commands/user/create.go' -> 'user'
	componentRegexp = regexp.MustCompile(`.*\/(.*)\/.*\.go`)
)

type Walk struct {
	projectDir                    string
	rootNamespace                 string
	componentsAndTheirConnections map[string][]string
}

func New(projectDir string, rootNamespace string) *Walk {
	return &Walk{
		projectDir:                    projectDir,
		rootNamespace:                 rootNamespace,
		componentsAndTheirConnections: map[string][]string{},
	}
}

func (w *Walk) FindComponents() error {
	return filepath.Walk(w.projectDir+"/"+w.rootNamespace, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isGoSourceFile(path) {
			regexpMatches := componentRegexp.FindStringSubmatch(path)

			if len(regexpMatches) == 2 {
				w.componentsAndTheirConnections[regexpMatches[1]] = []string{}
			}
		}

		return nil
	})
}

func (w *Walk) FindComponentConnections() error {
	return filepath.Walk(w.projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isGoSourceFile(path) {

		}

		return nil
	})
}

func (w *Walk) PrintComponents() {
	for component := range w.componentsAndTheirConnections {
		fmt.Println(component)
	}
}
