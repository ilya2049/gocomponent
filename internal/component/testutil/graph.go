package testutil

import "strings"

func BuildGraphString(imports ...string) string {
	sb := strings.Builder{}

	for _, imp := range imports {
		sb.WriteString(imp + "\n")
	}

	return sb.String()
}
