package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/artarts36/singlecli/color"
)

type helpCmd struct {
	Name          string
	Args          []*ArgDefinition
	Opts          []*OptDefinition
	UsageExamples []*UsageExample
}

func newHelpCmd(name string, args []*ArgDefinition, opts []*OptDefinition, usageExamples []*UsageExample) cmd {
	return (&helpCmd{
		Name:          name,
		Args:          args,
		Opts:          opts,
		UsageExamples: usageExamples,
	}).run
}

func (c *helpCmd) run(_ context.Context) error {
	signature := []string{
		c.Name,
	}

	if len(c.Args) > 0 {
		signature = append(signature, " ")

		for i, arg := range c.Args {
			signature = append(signature, arg.Name)

			if i != len(c.Args) {
				signature = append(signature, " ")
			}
		}
	}

	if len(c.Opts) > 0 {
		for i, opt := range c.Opts {
			if opt.WithValue {
				signature = append(signature, fmt.Sprintf("[--%s=<value>]", opt.Name))
			} else {
				signature = append(signature, fmt.Sprintf("[--%s]", opt.Name))
			}

			if i < len(c.Opts)-1 {
				signature = append(signature, " ")
			}
		}
	}

	fmt.Println(color.Yellow("Usage"))

	fmt.Printf("  %s\n", strings.Join(signature, ""))

	if len(c.Args) > 0 || len(c.Opts) > 0 {
		fmt.Println()
	}

	leftOffset := c.findLeftOffset()

	if len(c.Args) > 0 {
		fmt.Println(color.Yellow("Arguments"))

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

			fmt.Printf(
				"  %s%s%s%s%s\n",
				color.Green(arg.Name),
				strings.Repeat(" ", spaces),
				arg.Description,
				required,
				valuesEnum,
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

			fmt.Printf("  %s%s%s\n", color.Green(opt.Name), strings.Repeat(" ", spaces), opt.Description)
		}
	}

	if len(c.UsageExamples) > 0 {
		if len(c.Args) > 0 || len(c.Opts) > 0 {
			fmt.Println()
		}
		fmt.Println(color.Yellow("Usage examples"))

		for _, example := range c.UsageExamples {
			spaces := strings.Repeat(" ", leftOffset-len(example.Command))

			fmt.Printf(
				"  %s%s%s\n",
				color.Green(example.Command),
				spaces,
				example.Description,
			)
		}
	}

	return nil
}

func (c *helpCmd) findLeftOffset() int {
	const minOffset = 2

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

	return maxLen + minOffset
}
