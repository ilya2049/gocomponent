package component

type Component struct {
	namespace Namespace
	id        string
}

func New(namespace Namespace) *Component {
	return &Component{
		namespace: namespace,
		id:        "",
	}
}

func (c *Component) ExtendID() {
	c.id = c.namespace.ExtendComponentID(c.id)
}

func (c *Component) ID() string {
	return c.id
}
