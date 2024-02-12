package component

type Component struct {
	namespace    Namespace
	id           string
	isThirdParty bool
}

func New(namespace Namespace) *Component {
	return &Component{
		namespace:    namespace,
		id:           "",
		isThirdParty: false,
	}
}

func (c *Component) MarkAsThirdParty() {
	c.isThirdParty = true
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
