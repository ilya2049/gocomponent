package component

type Component struct {
	namespace   Namespace
	id          string
	isInProject bool
}

func New(namespace Namespace) *Component {
	return &Component{
		namespace:   namespace,
		id:          "",
		isInProject: false,
	}
}

func (c *Component) InProject() {
	c.isInProject = true
}

func (c *Component) IsInProject() bool {
	return c.isInProject
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
