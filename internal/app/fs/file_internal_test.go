package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isGoSourceFile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: ".golangci.yml is not a go file",
			path: ".golangci.yml",
			want: false,
		},
		{
			name: "summary.go is a go file",
			path: "summary.go",
			want: true,
		},
		{
			name: "internal/summary.go is a go file",
			path: "internal/summary.go",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isGoSourceFile(tt.path))
		})
	}
}
