package component_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"

	"github.com/stretchr/testify/assert"
)

func TestNamespace_ExcludeLastSection(t *testing.T) {
	tests := []struct {
		name string
		id   component.Namespace
		want string
	}{
		{
			name: "The last section is not empty",
			id:   "/internal/postgresql/connection",
			want: "/internal/postgresql/",
		},
		{
			name: "Namespace is a empty string",
			id:   "",
			want: "",
		},
		{
			name: "Namespace has only one section",
			id:   "connection",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.id.ExcludeLastSection())
		})
	}
}

func TestNamespace_LastSection(t *testing.T) {
	tests := []struct {
		name string
		id   component.Namespace
		want string
	}{
		{
			name: "The last section is not empty",
			id:   "/internal/postgresql/connection",
			want: "connection",
		},
		{
			name: "Namespace is a empty string",
			id:   "",
			want: "",
		},
		{
			name: "Namespace has only one section",
			id:   "connection",
			want: "connection",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.id.LastSection())
		})
	}
}

func TestNamespace_ExtendComponentID(t *testing.T) {
	tests := []struct {
		name    string
		id      component.Namespace
		shortID string
		want    string
	}{
		{
			name:    "There are not sections in a component id yet",
			id:      "/internal/postgresql/connection",
			shortID: "",
			want:    "connection",
		},
		{
			name:    "There is only one section in a component id",
			id:      "/internal/postgresql/connection",
			shortID: "connection",
			want:    "postgresql/connection",
		},
		{
			name:    "There are two sections in a component id",
			id:      "/internal/postgresql/connection",
			shortID: "postgresql/connection",
			want:    "internal/postgresql/connection",
		},
		{
			name:    "A component id cannot be extended",
			id:      "/internal/postgresql/connection",
			shortID: "internal/postgresql/connection",
			want:    "internal/postgresql/connection",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.id.ExtendComponentID(tt.shortID))
		})
	}
}
