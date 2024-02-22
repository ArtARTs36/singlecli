package cli

import (
	"context"
	"fmt"

	"github.com/artarts36/singlecli/color"
)

type versionCmd struct {
	*BuildInfo
}

func (c *versionCmd) Run(_ context.Context) error {
	fmt.Println(fmt.Sprintf(
		"%s version %s %s",
		color.Green(c.Name),
		color.Yellow(c.Version),
		c.BuildDate,
	))

	return nil
}
