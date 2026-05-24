package app

type closer struct {
	fns []func()
}

func (c *closer) add(fn func()) {
	c.fns = append(c.fns, fn)
}

func (c *closer) closeAll() {
	for i := len(c.fns) - 1; i >= 0; i-- {
		c.fns[i]()
	}
}
