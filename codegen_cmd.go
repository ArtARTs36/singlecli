package cli

import (
	"context"
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type CodegenFile struct {
	Name    string
	Content string
}

type githubAction struct {
	Name        string                                             `yaml:"name"`
	Description string                                             `yaml:"description"`
	Inputs      *orderedmap.OrderedMap[string, *githubActionInput] `yaml:"inputs"`
	Branding    struct {
		Icon  string `yaml:"icon"`
		Color string `yaml:"color,omitempty"`
	} `yaml:"branding,omitempty"`
	Runs struct {
		Using      string              `yaml:"using,omitempty"`
		Image      string              `yaml:"image,omitempty"`
		Entrypoint string              `yaml:"entrypoint,omitempty"`
		Args       []string            `yaml:"args,omitempty"`
		Steps      []*githubActionStep `yaml:"steps,omitempty"`
	} `yaml:"runs"`
}

type githubActionInput struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default,omitempty"`
}

type githubActionStep struct {
	Name             string                                      `yaml:"name"`
	ID               string                                      `yaml:"id,omitempty"`
	If               string                                      `yaml:"if,omitempty"`
	Run              string                                      `yaml:"run,omitempty"`
	Shell            string                                      `yaml:"shell,omitempty"`
	Uses             string                                      `yaml:"uses,omitempty"`
	With             map[string]string                           `yaml:"with,omitempty"`
	ContinueOnError  interface{}                                 `yaml:"continue-on-error,omitempty"`
	WorkingDirectory interface{}                                 `yaml:"working-directory,omitempty"`
	Env              *orderedmap.OrderedMap[string, interface{}] `yaml:"env,omitempty"`
}

type codegenGaCmd struct {
	app *App
}

func newCodegenGACmd(app *App) cmd {
	return (&codegenGaCmd{
		app: app,
	}).run
}

func (c *codegenGaCmd) run(_ context.Context) error {
	contentBytes, err := os.ReadFile("action.yaml")
	content := &githubAction{}
	if err == nil {
		err = yaml.Unmarshal(contentBytes, &content)
		if err != nil {
			return fmt.Errorf("failed to unmarshal action.yaml: %w", err)
		}
	}

	content = c.regenCode(content)

	contentYaml, err := yaml.Marshal(content)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	err = os.WriteFile("action.yaml", contentYaml, 0755)
	if err != nil {
		return fmt.Errorf("failed to save file action.yaml: %w", err)
	}

	return nil
}

func (c *codegenGaCmd) regenCode(
	content *githubAction,
) *githubAction {
	if len(c.app.Args) == 0 && len(c.app.Opts) == 0 {
		return content
	}

	inputs := orderedmap.New[string, *githubActionInput]()

	runArgs := make([]string, 0, len(c.app.Opts)+len(c.app.Args))

	for _, arg := range c.app.Args {
		defaultVal := ""

		if content.Inputs != nil {
			oldArg := content.Inputs.Value(arg.Name)
			if oldArg != nil && oldArg.Default != "" {
				defaultVal = oldArg.Default
			}
		}

		input := &githubActionInput{
			Description: arg.Description,
			Required:    arg.Required,
			Default:     defaultVal,
		}

		inputs.Set(arg.Name, input)

		runArgs = append(runArgs, fmt.Sprintf("${{ inputs.%s }}", arg.Name))
	}

	for _, opt := range c.app.Opts {
		defaultVal := ""

		if content.Inputs != nil {
			oldOpt := content.Inputs.Value(opt.Name)
			if oldOpt != nil && oldOpt.Default != "" {
				defaultVal = oldOpt.Default
			}
		}

		input := &githubActionInput{
			Description: opt.Description,
			Required:    false,
			Default:     defaultVal,
		}

		inputs.Set(opt.Name, input)

		if opt.WithValue {
			runArgs = append(runArgs, fmt.Sprintf("${{ inputs.%s != '' && format('--%s={0}', inputs.%s) || '' }}", opt.Name, opt.Name, opt.Name))
		} else {
			runArgs = append(runArgs, fmt.Sprintf("${{ inputs.%s != '' && '--%s' || '' }}", opt.Name, opt.Name))
		}
	}

	content.Inputs = inputs

	if content.Runs.Using == "docker" {
		content.Runs.Args = runArgs
	} else {
		for _, step := range content.Runs.Steps {
			if step.ID == "run-binary" {
				if step.Env == nil {
					step.Env = orderedmap.New[string, interface{}]()
				}

				args := runArgs
				for i, arg := range args {
					args[i] = fmt.Sprintf("'%s'", arg)
				}

				step.Env.Set("CMD_RUN_ARGS", &yaml.Node{
					Kind:  yaml.ScalarNode,
					Style: yaml.DoubleQuotedStyle,
					Value: strings.Join(args, " "),
				})
				continue
			}
		}
	}

	return content
}
