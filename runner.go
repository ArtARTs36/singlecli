package cli

import (
	"context"
	"fmt"
	"strings"
)

type actionRunner struct {
	Action         Action
	ArgDefinitions []*ArgDefinition
	OptDefinitions []*OptDefinition
	InputArgsList  []string
}

func (r *actionRunner) run(ctx context.Context) error {
	optDefMap := make(map[string]*OptDefinition)
	for _, def := range r.OptDefinitions {
		optDefMap[def.Name] = def
	}

	argMap := make(map[string]string)
	optMap := make(map[string]string)

	argumentsQueue := newArgQueue(r.ArgDefinitions)

	for i := 1; i < len(r.InputArgsList); i++ {
		val := r.InputArgsList[i]

		isOpt := strings.HasPrefix(val, "--")

		if isOpt {
			if argumentsQueue.valid() {
				required := argumentsQueue.firstRequired()
				if required != nil {
					return fmt.Errorf("must be set argument %q", required.Name)
				}

				argumentsQueue.clean()
			}

			optName := strings.SplitN(val, "--", 2)[1]
			_, exists := optDefMap[optName]
			if !exists {
				return fmt.Errorf("option %q unknown", optName)
			}

			optMap[optName] = "_"

			continue
		}

		def := argumentsQueue.pop()

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

	required := argumentsQueue.firstRequired()
	if required != nil {
		return fmt.Errorf("must be set argument %q", required.Name)
	}

	return r.Action(&Context{
		Context: ctx,
		Args:    argMap,
		Opts:    optMap,
	})
}
