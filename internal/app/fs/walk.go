package fs

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/ilya2049/gocomponent/internal/domain/component"
)

type Walk struct {
	fileReader     fileReader
	filePathWalker filePathWalker
	astFileParser  astFileParser
	projectDir     string
	project        *project
}

func NewWalk(
	projectDir string,
	fileReader fileReader,
	filePathWalker filePathWalker,
	astFileParser astFileParser,
) *Walk {
	if !strings.HasSuffix(projectDir, component.Slash) {
		projectDir += component.Slash
	}

	return &Walk{
		projectDir:     projectDir,
		project:        newProject(),
		fileReader:     fileReader,
		filePathWalker: filePathWalker,
		astFileParser:  astFileParser,
	}
}

func (w *Walk) ReadComponentGraph() (*component.Graph, error) {
	moduleName, err := w.readModuleName(w.projectDir)
	if err != nil {
		return nil, err
	}

	err = w.filePathWalker.Walk(w.projectDir, func(path string, info os.FileInfo, err error) error {
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

		if w.isRootNamespace(namespace) {
			namespace = component.NewNamespace(component.Slash)
		} else {
			namespace = component.Slash + namespace.TrimPrefix(w.projectDir)
		}

		aComponent := w.project.getOrAddComponent(namespace)

		packageImports, err := w.parseImportsOfGoFile(namespace, moduleName, path)
		if err != nil {
			return err
		}

		aPackage := newProjectPackage(aComponent, packageImports)

		w.addPackageInProject(namespace, aPackage)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return w.project.createComponentGraph(), nil
}

func (w *Walk) isRootNamespace(namespace component.Namespace) bool {
	return namespace+component.Slash == component.Namespace(w.projectDir)
}

func (w *Walk) addPackageInProject(namespace component.Namespace, newPackage *projectPackage) {
	existingPackage, ok := w.project.findPackage(namespace)
	if ok {
		existingPackage.join(newPackage)

		return
	}

	w.project.addPackage(namespace, newPackage)
}

func (w *Walk) parseImportsOfGoFile(
	currentNamespace component.Namespace,
	moduleName string,
	goFileName string,
) (map[component.Namespace]*component.Component, error) {
	file, err := w.astFileParser.ParseFile(token.NewFileSet(), goFileName, nil, parser.Mode(0))
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	imports := make(map[component.Namespace]*component.Component)

	for _, fileImport := range file.Imports {
		namespace := component.NewNamespace(fileImport.Path.Value[1 : len(fileImport.Path.Value)-1])

		moduleNameWithSectionSeparator := moduleName + component.Slash

		if namespace.HasPrefix(moduleNameWithSectionSeparator) {
			namespace = namespace.TrimPrefix(moduleName)
		}

		if namespace == currentNamespace {
			continue
		}

		component := w.project.getOrAddComponent(namespace)

		imports[namespace] = component
	}

	return imports, nil
}

var ErrFirstLineOfGoModShouldIncludeExactlyTwoPArts = errors.New(
	"the first line of go mod file parts should includes exactly two parts",
)

func (w *Walk) readModuleName(projectDir string) (string, error) {
	goModFileContents, err := w.fileReader.ReadFile(projectDir + "go.mod")
	if err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}

	var firstLineOfGoModFile = []byte{}

	for _, b := range goModFileContents {
		if b == '\n' {
			break
		} else {
			firstLineOfGoModFile = append(firstLineOfGoModFile, b)
		}
	}

	firstLineOfGoModFileParts := strings.Split(string(firstLineOfGoModFile), " ")
	if len(firstLineOfGoModFileParts) != 2 {
		return "", ErrFirstLineOfGoModShouldIncludeExactlyTwoPArts
	}

	return firstLineOfGoModFileParts[1], nil
}
