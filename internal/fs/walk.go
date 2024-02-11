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
	packages          map[component.Namespace]*component.Package
	isDebug           bool
}

func NewWalk(projectDir string, isDebug bool) *Walk {
	return &Walk{
		projectDir:        projectDir,
		componentRegistry: component.NewRegistry(isDebug),
		packages:          make(map[component.Namespace]*component.Package),
		isDebug:           isDebug,
	}
}

func (w *Walk) FindComponentsAndImports() error {
	moduleName, err := readModuleName(w.projectDir)
	if err != nil {
		return err
	}

	if w.isDebug {
		fmt.Println("DEBUG", "module", moduleName)
		fmt.Println("module", moduleName)
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
	existingPackage, ok := w.packages[namespace]
	if ok {
		existingPackage.Join(newPackage)

		return
	}

	w.packages[namespace] = newPackage
}

func (w *Walk) ConvertComponentsAndImportsToDotGraphDotGraph() string {
	if w.isDebug {
		fmt.Println("DEBUG", "dot graph")
	}

	sb := strings.Builder{}

	sb.WriteString("digraph {\n")

	for _, p := range w.packages {
		sb.WriteString(`"` + p.ID() + `"` + "\n")

		if w.isDebug {
			fmt.Println("package ns:", p.Namespace())
		}

		for _, importedComponent := range p.Imports() {
			if w.isDebug {
				fmt.Println("import ns:", importedComponent.Namespace())
			}

			sb.WriteString(`"` + p.ID() + `" -> "` + importedComponent.ID() + `"` + "\n")
		}
	}

	sb.WriteString("}")

	return sb.String()
}
