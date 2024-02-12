package fs

import (
	"fmt"
	"go/parser"
	"go/token"
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
	if !strings.HasSuffix(projectDir, "/") {
		projectDir += "/"
	}

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

	err = filepath.Walk(w.projectDir, func(path string, info os.FileInfo, err error) error {
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

		namespace = namespace.TrimPrefix(w.projectDir)
		aComponent := w.componentRegistry.GetOrAddComponent(namespace)

		packageImports, err := w.parseImportsOfGoFile(namespace, moduleName, path)
		if err != nil {
			return err
		}

		aPackage := component.NewPackage(aComponent, packageImports)

		w.addPackageInProject(namespace, aPackage)

		return nil
	})

	if err != nil {
		return err
	}

	w.componentRegistry.MakeUniqueComponentIDs()

	return nil
}

func (w *Walk) addPackageInProject(namespace component.Namespace, newPackage *component.Package) {
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

func (w *Walk) parseImportsOfGoFile(
	currentNamespace component.Namespace,
	moduleName string,
	goFileName string,
) (map[component.Namespace]*component.Component, error) {
	file, err := parser.ParseFile(token.NewFileSet(), goFileName, nil, parser.Mode(0))
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	imports := make(map[component.Namespace]*component.Component)

	for _, fileImport := range file.Imports {
		namespace := component.NewNamespace(fileImport.Path.Value[1 : len(fileImport.Path.Value)-1])

		var isComponentInProject bool

		moduleNameWithSlash := moduleName + "/"

		if namespace.HasPrefix(moduleNameWithSlash) {
			namespace = namespace.TrimPrefix(moduleNameWithSlash)
			isComponentInProject = true
		}

		if namespace == currentNamespace {
			continue
		}

		component := w.componentRegistry.GetOrAddComponent(namespace)
		if !isComponentInProject {
			component.MarkAsThirdParty()
		}

		imports[namespace] = component
	}

	return imports, nil
}
