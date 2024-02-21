package fs

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/project"
)

type Walk struct {
	projectDir string
	project    *project.Project
}

func NewWalk(projectDir string, prj *project.Project) *Walk {
	if !strings.HasSuffix(projectDir, "/") {
		projectDir += "/"
	}

	return &Walk{
		projectDir: projectDir,
		project:    prj,
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

		namespace = "/" + namespace.TrimPrefix(w.projectDir)

		if namespace+"/" == component.Namespace(w.projectDir) {
			namespace = component.NewNamespace(moduleName)
		}

		aComponent := w.project.GetOrAddComponent(namespace)

		packageImports, err := w.parseImportsOfGoFile(namespace, moduleName, path)
		if err != nil {
			return err
		}

		aPackage := project.NewPackage(aComponent, packageImports)

		w.addPackageInProject(namespace, aPackage)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (w *Walk) addPackageInProject(namespace component.Namespace, newPackage *project.Package) {
	existingPackage, ok := w.project.FindPackage(namespace)
	if ok {
		existingPackage.Join(newPackage)

		return
	}

	w.project.AddPackage(namespace, newPackage)
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
			namespace = namespace.TrimPrefix(moduleName)
			isComponentInProject = true
		}

		if namespace == currentNamespace {
			continue
		}

		component := w.project.GetOrAddComponent(namespace)
		if !isComponentInProject {
			component.MarkAsThirdParty()
		}

		imports[namespace] = component
	}

	return imports, nil
}
