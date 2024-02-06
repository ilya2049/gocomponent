package component_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"

	"github.com/stretchr/testify/assert"
)

func TestID_Namespace(t *testing.T) {
	tests := []struct {
		name string
		id   component.ID
		want string
	}{
		{
			name: "The last part of id is a component name, the rest is a namespace",
			id:   "/internal/postgresql/connection",
			want: "/internal/postgresql/",
		},
		{
			name: "Id is a empty string",
			id:   "",
			want: "",
		},
		{
			name: "Id has no a namespace",
			id:   "connection",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.id.Namespace())
		})
	}
}

func TestID_ComponentName(t *testing.T) {
	tests := []struct {
		name string
		id   component.ID
		want string
	}{
		{
			name: "A component name is the last part of id",
			id:   "/internal/postgresql/connection",
			want: "connection",
		},
		{
			name: "Id is a empty string",
			id:   "",
			want: "",
		},
		{
			name: "Id has no a namespace",
			id:   "connection",
			want: "connection",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.id.ComponentName())
		})
	}
}

func TestID_ExtendShortID(t *testing.T) {
	tests := []struct {
		name    string
		id      component.ID
		shortID string
		want    string
	}{
		{
			name:    "A short id is empty",
			id:      "/internal/postgresql/connection",
			shortID: "",
			want:    "connection",
		},
		{
			name:    "A short id is a component name",
			id:      "/internal/postgresql/connection",
			shortID: "connection",
			want:    "postgresql/connection",
		},
		{
			name:    "A short id has already extended",
			id:      "/internal/postgresql/connection",
			shortID: "postgresql/connection",
			want:    "internal/postgresql/connection",
		},
		{
			name:    "A short id cannot be extended",
			id:      "/internal/postgresql/connection",
			shortID: "internal/postgresql/connection",
			want:    "internal/postgresql/connection",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.id.ExtendShortID(tt.shortID))
		})
	}
}
