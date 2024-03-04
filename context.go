package cli

import "context"

type Context struct {
	Context context.Context
	Args    map[string]string
	Opts    map[string]string
}

func (c *Context) GetArg(name string) string {
	return c.Args[name]
}

func (c *Context) HasOpt(name string) bool {
	_, exists := c.Opts[name]

	return exists
}
