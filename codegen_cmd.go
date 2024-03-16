package cli

import (
	"context"
	"fmt"
	"os"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"
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
		Using      string   `yaml:"using"`
		Image      string   `yaml:"image"`
		Entrypoint string   `yaml:"entrypoint,omitempty"`
		Args       []string `yaml:"args,omitempty"`
	} `yaml:"runs"`
}

type githubActionInput struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default,omitempty"`
}

type codegenCmd struct {
	app *App
}

func newCodegenGACmd(app *App) cmd {
	return (&codegenCmd{
		app: app,
	}).run
}

func (c *codegenCmd) run(_ context.Context) error {
	contentBytes, err := os.ReadFile("action.yaml")
	var content *githubAction
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

func (c *codegenCmd) regenCode(
	content *githubAction,
) *githubAction {
	if len(c.app.Args) == 0 && len(c.app.Opts) == 0 {
		return content
	}

	inputs := orderedmap.New[string, *githubActionInput]()

	runArgs := make([]string, 0, len(c.app.Opts)+len(c.app.Args))

	for _, arg := range c.app.Args {
		input := &githubActionInput{
			Description: arg.Description,
			Required:    arg.Required,
		}

		inputs.Set(arg.Name, input)

		runArgs = append(runArgs, fmt.Sprintf("${{ inputs.%s }}", arg.Name))
	}

	for _, opt := range c.app.Opts {
		input := &githubActionInput{
			Description: opt.Description,
			Required:    false,
			Default:     "",
		}

		inputs.Set(opt.Name, input)

		if opt.WithValue {
			runArgs = append(runArgs, fmt.Sprintf("${{ inputs.%s != '' && format('--%s={0}', inputs.%s) || '' }}", opt.Name, opt.Name, opt.Name))
		} else {
			runArgs = append(runArgs, fmt.Sprintf("${{ inputs.%s != '' && '--%s' || '' }}", opt.Name, opt.Name))
		}
	}

	content.Inputs = inputs
	content.Runs.Args = runArgs

	return content
}
