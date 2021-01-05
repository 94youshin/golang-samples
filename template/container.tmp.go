package PACKAGE_NAME

type GENERIC_NAMEContainer struct {
	s []GENERAC_TYPE
}

func NewGENERAC_NAMEContainer() * GENERAC_NAMEContainer {
	return &GENERIC_NAMEContainer{s: []gENERIC_TYPE{}}
}

func (c *GENERIC_NAMEContainer) Put (val GENERIC_TYPE) {
	c.s = append(c.s, val)
}

func (c *GENERIC_NAMEContainer) Get() GENERIC_TYPE {
	r := c.s[0]
	c.s = c.s[1:]
	return r
}