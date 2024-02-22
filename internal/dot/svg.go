package dot

import (
	"strings"

	"github.com/ilya2049/gocomponent/internal/component"
)

func MapComponentAndNamespace(components component.Components, svgBytes []byte) []byte {
	svgString := string(svgBytes)

	for _, aComponent := range components {
		svgString = strings.Replace(svgString,
			"<title>"+aComponent.ID()+"</title>",
			"<title>"+aComponent.Namespace()+"</title>",
			1,
		)
	}

	return []byte(svgString)
}
