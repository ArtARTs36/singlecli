package cli

import (
	"context"
	"fmt"
	"os"
)

type App struct {
	BuildInfo     *BuildInfo
	Args          []*ArgDefinition
	Action        Action
	UsageExamples []*UsageExample
}

type BuildInfo struct {
	Name      string
	Version   string
	BuildDate string
}

type Action func(ctx *Context) error

type cmd func(ctx context.Context) error

func (a *App) RunWithGlobalArgs(ctx context.Context) {
	a.Run(ctx, os.Args[1:])
}

func (a *App) Run(ctx context.Context, args []string) {
	c := a.findCmd(args)

	err := c(ctx)
	if err != nil {
		fmt.Printf("action failed: %s\n", err)

		os.Exit(1)
	}

	os.Exit(0)
}

func (a *App) findCmd(args []string) cmd {
	if len(args) == 0 {
		return newHelpCmd(a.Args, a.UsageExamples)
	}

	for _, arg := range args {
		if arg == "--version" {
			return (&versionCmd{
				BuildInfo: a.BuildInfo,
			}).Run
		}

		if arg == "--help" {
			return newHelpCmd(a.Args, a.UsageExamples)
		}
	}

	return (&actionRunner{
		Action:         a.Action,
		ArgDefinitions: a.Args,
		InputArgsList:  args,
	}).run
}
