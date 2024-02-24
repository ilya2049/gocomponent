package component

import "strings"

type Component struct {
	namespace Namespace
	id        string
	color     Color
}

func New(namespace Namespace) *Component {
	return &Component{
		namespace: namespace,
		id:        "",
		color:     "",
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

type Components []*Component
