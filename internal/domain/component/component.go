package component

import "strings"

type Component struct {
	namespace      Namespace
	id             string
	color          Color
	sizeBytes      int
	normalizedSize float64
}

func New(namespace Namespace) *Component {
	return &Component{
		namespace:      namespace,
		id:             "",
		color:          "",
		sizeBytes:      0,
		normalizedSize: 0,
	}
}

func (c *Component) IsThirdParty() bool {
	return !strings.HasPrefix(c.namespace.String(), Slash)
}

func (c *Component) ExtendID() {
	c.id = c.namespace.ExtendComponentID(c.id)
}

func (c *Component) ID() string {
	return c.id
}

func (c *Component) Namespace() Namespace {
	return c.namespace
}

func (c *Component) Colorize(color Color) {
	c.color = color
}

func (c *Component) Color() Color {
	return c.color
}

func (c *Component) HasColor() bool {
	return c.color != ""
}

func (c *Component) AddBytesInSize(bytes int) {
	if bytes < 0 {
		return
	}

	c.sizeBytes += bytes
}

func (c *Component) SizeBytes() int {
	return c.sizeBytes
}

func (c *Component) NormalizeSize(normalizedSize float64) {
	c.normalizedSize = normalizedSize
}

func (c *Component) NormalizedSize() float64 {
	return c.normalizedSize
}

type Components []*Component
