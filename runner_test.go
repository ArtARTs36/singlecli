package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type actionRunnerParseArgsExpectedResult struct {
	Args  map[string]string
	Opts  map[string]string
	Error error
}

func Test_ActionRunner_ParseArgs(t *testing.T) {
	tCases := []struct {
		Name string

		ArgDefinitions []*ArgDefinition
		OptDefinitions []*OptDefinition
		InputArgsList  []string

		ExpectedResult actionRunnerParseArgsExpectedResult
	}{
		{
			Name: "with empty args",

			ArgDefinitions: []*ArgDefinition{},
			OptDefinitions: []*OptDefinition{},
			InputArgsList:  []string{},
			ExpectedResult: actionRunnerParseArgsExpectedResult{
				Args: map[string]string{},
				Opts: map[string]string{},
			},
		},
		{
			Name: "with args",

			ArgDefinitions: []*ArgDefinition{
				{
					Name: "name",
				},
				{
					Name: "age",
				},
			},
			OptDefinitions: []*OptDefinition{},
			InputArgsList: []string{
				"",
				"artem",
				"24",
			},

			ExpectedResult: actionRunnerParseArgsExpectedResult{
				Args: map[string]string{
					"name": "artem",
					"age":  "24",
				},
				Opts: map[string]string{},
			},
		},
		{
			Name: "with args and opts",

			ArgDefinitions: []*ArgDefinition{
				{
					Name: "name",
				},
				{
					Name: "age",
				},
			},
			OptDefinitions: []*OptDefinition{
				{
					Name: "lang",
				},
			},
			InputArgsList: []string{
				"",
				"artem",
				"24",
				"--lang",
			},

			ExpectedResult: actionRunnerParseArgsExpectedResult{
				Args: map[string]string{
					"name": "artem",
					"age":  "24",
				},
				Opts: map[string]string{
					"lang": "",
				},
			},
		},
		{
			Name: "with args and opts with values",

			ArgDefinitions: []*ArgDefinition{
				{
					Name: "name",
				},
				{
					Name: "age",
				},
			},
			OptDefinitions: []*OptDefinition{
				{
					Name: "lang",
				},
				{
					Name: "weight",
				},
			},
			InputArgsList: []string{
				"",
				"artem",
				"24",
				"--lang",
				"--weight=80",
			},

			ExpectedResult: actionRunnerParseArgsExpectedResult{
				Args: map[string]string{
					"name": "artem",
					"age":  "24",
				},
				Opts: map[string]string{
					"lang":   "",
					"weight": "80",
				},
			},
		},
	}

	for _, tCase := range tCases {
		t.Run(tCase.Name, func(t *testing.T) {
			runner := &actionRunner{
				ArgDefinitions: tCase.ArgDefinitions,
				OptDefinitions: tCase.OptDefinitions,
				InputArgsList:  tCase.InputArgsList,
			}

			args, opts, err := runner.parseArgs()

			assert.Equal(
				t,
				tCase.ExpectedResult,

				actionRunnerParseArgsExpectedResult{
					Args:  args,
					Opts:  opts,
					Error: err,
				},
			)
		})
	}
}
