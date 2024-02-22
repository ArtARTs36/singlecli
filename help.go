package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/artarts36/singlecli/color"
)

type helpCmd struct {
	Args          []*ArgDefinition
	UsageExamples []*UsageExample
}

func newHelpCmd(args []*ArgDefinition, usageExamples []*UsageExample) cmd {
	return (&helpCmd{
		Args:          args,
		UsageExamples: usageExamples,
	}).run
}

func (c *helpCmd) run(_ context.Context) error {
	fmt.Println(color.Yellow("Arguments"))

	leftOffset := c.findLeftOffset()

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

	if len(c.UsageExamples) > 0 {
		fmt.Println()
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

	return maxLen + 2
}
