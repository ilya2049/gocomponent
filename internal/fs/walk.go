package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilya2049/gocomponent/internal/regexp"
)

type Walk struct {
	projectDir                 string
	rootNamespace              string
	componentsAndImports       map[string]map[string]struct{}
	isComponentsHaveDoubleName bool
}

func NewWalk(projectDir string, rootNamespace string) *Walk {
	return &Walk{
		projectDir:           projectDir,
		rootNamespace:        rootNamespace,
		componentsAndImports: map[string]map[string]struct{}{},
	}
}

func (w *Walk) ComponentsHaveDoubleName() {
	w.isComponentsHaveDoubleName = true
}

func (w *Walk) FindComponentsAndImports() error {
	componentImportRegexp, err := regexp.CompileComponentImportRegexp(w.rootNamespace, w.isComponentsHaveDoubleName)
	if err != nil {
		return fmt.Errorf("compile component import regexp: %w", err)
	}

	return filepath.Walk(w.projectDir+"/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk err: %w", err)
		}

		if isGoSourceFile(path) {
			if component := regexp.FindComponent(path, w.isComponentsHaveDoubleName); component != "" {
				w.saveComponentWithoutImport(component)

				goFileContents, err := readFile(path)
				if err != nil {
					return fmt.Errorf("read .go file: %w", err)
				}

				componentImports := regexp.FindComponentImports(componentImportRegexp, goFileContents)

				for _, componentImport := range componentImports {
					w.saveComponentImport(component, componentImport)
				}
			}
		}

		return nil
	})
}

func (w *Walk) saveComponentWithoutImport(component string) {
	if _, ok := w.componentsAndImports[component]; !ok {
		w.componentsAndImports[component] = map[string]struct{}{}
	}
}

func (w *Walk) saveComponentImport(component, componentImport string) {
	if component == componentImport {
		return
	}

	w.componentsAndImports[component][componentImport] = struct{}{}
}

func (w *Walk) ConvertComponentsAndImportsToDotGraphDotGraph() string {
	sb := strings.Builder{}

	sb.WriteString("digraph G {\n")

	for component, componentImports := range w.componentsAndImports {
		sb.WriteString(`"` + component + `"` + "\n")

		if len(componentImports) > 0 {
			for componentImport := range componentImports {
				sb.WriteString(`"` + component + `" -> "` + componentImport + `"` + "\n")
			}
		}
	}

	sb.WriteString("}")

	return sb.String()
}
