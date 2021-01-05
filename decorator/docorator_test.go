package decorator

import (
	"testing"
)

func TestWarpAddDecorator(t *testing.T) {
	var c Component = &ConcreteComponent{}
	c = WarpAddDecorator(c, 10)
	c = WarpMulDecorator(c, 8)
	res := c.Calc()
	t.Log(res)
}
