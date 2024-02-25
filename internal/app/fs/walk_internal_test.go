package fs

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/pkg/fs"
	"github.com/ilya2049/gocomponent/internal/pkg/sbuilder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalk_readModuleName_parseGoModContents(t *testing.T) {
	// Given
	goModFileBytes := []byte(sbuilder.BuildMultilineString(
		"module github.com/ilya2049/gocomponent",
		"",
		"go 1.21",
		"",
		"require (",
		"    github.com/BurntSushi/toml v1.3.2",
		")",
	))

	walk := &Walk{
		fileReader: fs.NewFileReaderStub(goModFileBytes, nil),
	}

	// When
	moduleName, err := walk.readModuleName("")
	require.NoError(t, err)

	// Then
	assert.Equal(t, "github.com/ilya2049/gocomponent", moduleName)
}
