package component

type Component struct {
	namespace    Namespace
	id           string
	isThirdParty bool
	color        Color
}

func New(namespace Namespace) *Component {
	return &Component{
		namespace:    namespace,
		id:           "",
		isThirdParty: false,
		color:        "",
	}
}

func (c *Component) MarkAsThirdParty() {
	c.isThirdParty = true
}

func (c *Component) IsThirdParty() bool {
	return c.isThirdParty
}

func (c *Component) ExtendID() {
	c.id = c.namespace.ExtendComponentID(c.id)
}

func (c *Component) ID() string {
	return c.id
}

func (c *Component) Namespace() string {
	return string(c.namespace)
}

func (c *Component) Colorize(color Color) {
	c.color = color
}

func (c *Component) HasColor() bool {
	return c.color != ""
}

func (c *Component) Color() Color {
	return c.color
}

type Components []*Component
