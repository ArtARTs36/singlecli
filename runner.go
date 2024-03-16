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
	argMap, optMap, err := r.parseArgs()
	if err != nil {
		return err
	}

	return r.Action(&Context{
		Context: ctx,
		Output:  &output{},
		Args:    argMap,
		Opts:    optMap,
	})
}

func (r *actionRunner) parseArgs() (map[string]string, map[string]string, error) {
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
					return nil, nil, fmt.Errorf("must be set argument %q", required.Name)
				}

				argumentsQueue.clean()
			}

			optParts := strings.SplitN(
				strings.SplitN(val, "--", 2)[1],
				"=",
				2,
			)

			optName := optParts[0]
			optValue := ""

			if len(optParts) == 2 {
				optValue = optParts[1]
			}

			_, exists := optDefMap[optName]
			if !exists {
				return nil, nil, fmt.Errorf("option %q unknown", optName)
			}

			optMap[optName] = optValue

			continue
		}

		def := argumentsQueue.pop()

		if def == nil {
			return nil, nil, fmt.Errorf("unknown arg with value %q", val)
		}

		if def.HasValuesEnum() && !def.ValueInputs(val) {
			return nil, nil, fmt.Errorf(
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
		return nil, nil, fmt.Errorf("must be set argument %q", required.Name)
	}

	return argMap, optMap, nil
}
