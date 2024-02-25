package sbuilder

import "strings"

func BuildMultilineString(lines ...string) string {
	sb := strings.Builder{}

	for i, line := range lines {
		sb.WriteString(line)

		if i < len(lines)-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}
