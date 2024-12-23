package cli

import (
	"context"
	"github.com/artarts36/singlecli/color"
	"os"
)

type App struct {
	BuildInfo     *BuildInfo
	Args          []*ArgDefinition
	Opts          []*OptDefinition
	Action        Action
	UsageExamples []*UsageExample
}

type BuildInfo struct {
	Name        string
	Description string
	Version     string
	BuildDate   string
}

type Action func(ctx *Context) error

type cmd func(ctx context.Context) error

func (a *App) RunWithGlobalArgs(ctx context.Context) {
	a.Run(ctx, os.Args)
}

func (a *App) Run(ctx context.Context, args []string) {
	c := a.findCmd(args)

	err := c(ctx)
	if err != nil {
		output{}.PrintColoredBlock(color.Red, err.Error())

		os.Exit(1)
	}

	os.Exit(0)
}

func (a *App) findCmd(args []string) cmd {
	if len(args) == 1 && (len(a.Args) > 0 && a.Args[0].Required) {
		return newHelpCmd(a.BuildInfo.Name, a.Args, a.Opts, a.UsageExamples)
	}

	for _, arg := range args {
		if arg == "--version" {
			return (&versionCmd{
				BuildInfo: a.BuildInfo,
			}).Run
		}

		if arg == "--singlecli-codegen-ga" {
			return newCodegenGACmd(a)
		}

		if arg == "--help" {
			return newHelpCmd(a.BuildInfo.Name, a.Args, a.Opts, a.UsageExamples)
		}
	}

	return (&actionRunner{
		Action:         a.Action,
		ArgDefinitions: a.Args,
		OptDefinitions: a.Opts,
		InputArgsList:  args,
	}).run
}
