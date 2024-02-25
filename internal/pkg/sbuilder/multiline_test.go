package sbuilder_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/pkg/sbuilder"
	"github.com/stretchr/testify/assert"
)

func TestMakeMultilineString(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{
			name: "There is one line",
			lines: []string{
				"line1",
			},
			want: "line1",
		},
		{
			name: "There are two lines",
			lines: []string{
				"line1",
				"line2",
			},
			want: "line1\nline2",
		},
		{
			name: "There are two lines and a empty line between them",
			lines: []string{
				"line1",
				"",
				"line2",
			},
			want: "line1\n\nline2",
		},
		{
			name: "The lines have spaces",
			lines: []string{
				"line 1",
				"line 2",
			},
			want: "line 1\nline 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, sbuilder.BuildMultilineString(tt.lines...))
		})
	}
}
