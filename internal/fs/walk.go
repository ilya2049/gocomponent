package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilya2049/gocomponent/internal/component"
)

type Walk struct {
	projectDir        string
	componentRegistry *component.Registry
	project           *component.Project
}

func NewWalk(projectDir string) *Walk {
	return &Walk{
		projectDir:        projectDir,
		componentRegistry: component.NewRegistry(),
		project:           component.NewProject(),
	}
}

func (w *Walk) FindComponentsAndImports() error {
	moduleName, err := readModuleName(w.projectDir)
	if err != nil {
		return err
	}

	err = filepath.Walk(w.projectDir+"/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk err: %w", err)
		}

		if !isGoSourceFile(path) {
			return nil
		}

		namespace, ok := findNamespaceInPath(path)
		if !ok {
			return nil
		}

		namespace = namespace.TrimPrefix(w.projectDir + "/")

		c := w.componentRegistry.GetOrAddComponent(namespace)

		p := component.NewPackage(c)
		p.ParseImportsOfGoFile(moduleName, path, w.componentRegistry)

		w.addPackage(namespace, p)

		return nil
	})

	if err != nil {
		return err
	}

	w.componentRegistry.MakeUniqueComponentIDs()

	return nil
}

func (w *Walk) addPackage(namespace component.Namespace, newPackage *component.Package) {
	existingPackage, ok := w.project.FindPackage(namespace)
	if ok {
		existingPackage.Join(newPackage)

		return
	}

	w.project.AddPackage(namespace, newPackage)
}

func (w *Walk) ConvertComponentsAndImportsToDotGraphDotGraph(showThirdPartyImports bool) string {
	sb := strings.Builder{}

	sb.WriteString("digraph {\n")

	for _, p := range w.project.Packages() {
		sb.WriteString(`"` + p.ID() + `"` + "\n")

		for _, importedComponent := range p.Imports() {

			if showThirdPartyImports {
				sb.WriteString(`"` + p.ID() + `" -> "` + importedComponent.ID() + `"` + "\n")
			} else {
				if !importedComponent.IsThirdParty() {
					sb.WriteString(`"` + p.ID() + `" -> "` + importedComponent.ID() + `"` + "\n")
				}
			}
		}
	}

	sb.WriteString("}")

	return sb.String()
}
