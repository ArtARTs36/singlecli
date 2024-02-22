package cli

import "context"

type Context struct {
	Context context.Context
	Args    map[string]string
}

func (c *Context) GetArg(name string) string {
	return c.Args[name]
}
