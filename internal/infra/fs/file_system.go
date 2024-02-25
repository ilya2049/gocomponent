package fs

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/ilya2049/gocomponent/internal/app/fs"
	"github.com/ilya2049/gocomponent/internal/domain/component"
	"github.com/ilya2049/gocomponent/internal/infra/config"
)

type ComponentGraphReader struct {
}

func (r *ComponentGraphReader) ReadComponentGraph() (*component.Graph, error) {
	configReader := config.NewReader(&fileReader{})

	conf, err := configReader.ReadConfig()
	if err != nil {
		return nil, err
	}

	fsWalk := fs.NewWalk(
		conf.ProjectDirectory,
		&fileReader{},
		&filePathWalker{},
		&astFileParser{},
	)

	initialComponentGraph, err := fsWalk.ReadComponentGraph()
	if err != nil {
		return nil, err
	}

	componentGraph, err := component.ApplyGraphConfig(conf.ToComponentGraphConfig(), initialComponentGraph)
	if err != nil {
		return nil, err
	}

	return componentGraph, nil
}

type fileReader struct {
}

func (r *fileReader) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

type filePathWalker struct {
}

func (w *filePathWalker) Walk(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(root, fn)
}

type astFileParser struct {
}

func (p *astFileParser) ParseFile(fset *token.FileSet, filename string, src any, mode parser.Mode) (f *ast.File, err error) {
	return parser.ParseFile(fset, filename, src, mode)
}
