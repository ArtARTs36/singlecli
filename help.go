package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/artarts36/singlecli/color"
)

type helpCmd struct {
	Args          []*ArgDefinition
	Opts          []*OptDefinition
	UsageExamples []*UsageExample
}

func newHelpCmd(args []*ArgDefinition, opts []*OptDefinition, usageExamples []*UsageExample) cmd {
	return (&helpCmd{
		Args:          args,
		Opts:          opts,
		UsageExamples: usageExamples,
	}).run
}

func (c *helpCmd) run(_ context.Context) error {
	fmt.Println(color.Yellow("Arguments"))

	leftOffset := c.findLeftOffset()

	if len(c.Args) > 0 {
		for _, arg := range c.Args {
			spaces := leftOffset - len(arg.Name)

			required := ""
			if arg.Required {
				required = ", required"
			}

			valuesEnum := ""
			if len(arg.ValuesEnum) > 0 {
				for i, v := range arg.ValuesEnum {
					valuesEnum += v

					if i < len(arg.ValuesEnum)-1 {
						valuesEnum += ", "
					}
				}

				valuesEnum = fmt.Sprintf(", available values: [%s]", valuesEnum)
			}

			fmt.Println(
				fmt.Sprintf(
					"  %s%s%s%s%s",
					color.Green(arg.Name),
					strings.Repeat(" ", spaces),
					arg.Description,
					required,
					valuesEnum,
				),
			)
		}
	}

	if len(c.Opts) > 0 {
		if len(c.Args) > 0 {
			fmt.Println()
		}
		fmt.Println(color.Yellow("Options"))

		for _, opt := range c.Opts {
			spaces := leftOffset - len(opt.Name)

			fmt.Println(
				fmt.Sprintf(
					"  %s%s%s",
					color.Green(opt.Name),
					strings.Repeat(" ", spaces),
					opt.Description,
				),
			)
		}
	}

	if len(c.UsageExamples) > 0 {
		if len(c.Args) > 0 || len(c.Opts) > 0 {
			fmt.Println()
		}
		fmt.Println(color.Yellow("Usage examples"))

		for _, example := range c.UsageExamples {
			spaces := strings.Repeat(" ", leftOffset-len(example.Command))

			fmt.Println(fmt.Sprintf(
				"  %s%s%s",
				color.Green(example.Command),
				spaces,
				example.Description,
			))
		}
	}

	return nil
}

func (c *helpCmd) findLeftOffset() int {
	maxLen := 0

	for _, arg := range c.Args {
		if maxLen < len(arg.Name) {
			maxLen = len(arg.Name)
		}
	}

	if len(c.UsageExamples) > 0 {
		for _, example := range c.UsageExamples {
			if len(example.Command) > maxLen {
				maxLen = len(example.Command)
			}
		}
	}

	for _, opt := range c.Opts {
		if maxLen < len(opt.Name) {
			maxLen = len(opt.Name)
		}
	}

	return maxLen + 2
}
