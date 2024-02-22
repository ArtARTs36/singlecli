package cli

import (
	"context"
	"fmt"
	"strings"
)

type actionRunner struct {
	Action         Action
	ArgDefinitions []*ArgDefinition
	InputArgsList  []string
}

func (r *actionRunner) run(ctx context.Context) error {
	argMap := make(map[string]string)

	for i, def := range r.ArgDefinitions {
		i += 1

		argExists := i <= len(r.InputArgsList)-1
		if !argExists || r.InputArgsList[i] == "" {
			if def.Required {
				return fmt.Errorf("must be set argument %q", def.Name)
			}

			continue
		}

		val := r.InputArgsList[i]

		if def.HasValuesEnum() && !def.ValueInputs(val) {
			return fmt.Errorf(
				"value %q for argument %q invalid. Available values: [%s]",
				val,
				def.Name,
				strings.Join(def.ValuesEnum, ", "),
			)
		}

		argMap[def.Name] = val
	}

	return r.Action(&Context{
		Context: ctx,
		Args:    argMap,
	})
}
