package dot

import (
	"strings"

	"github.com/ilya2049/gocomponent/internal/component"
)

func MapComponentAndNamespaceInSVG(components component.Components, svgBytes []byte) []byte {
	svgString := string(svgBytes)

	for _, aComponent := range components {
		svgString = strings.Replace(svgString,
			"<title>"+aComponent.ID()+"</title>",
			"<title>"+aComponent.Namespace().String()+"</title>",
			1,
		)
	}

	return []byte(svgString)
}
