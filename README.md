# singlecli

singlecli - package for creating single command console application

## App definition

It's simple app for adding numbers

```go
package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	cli "github.com/artarts36/singlecli"
)

func main() {
	application := cli.App{
		BuildInfo: &cli.BuildInfo{
			Name:      "number-adder",
			Version:   "0.1.0",
			BuildDate: time.Now().String(),
		},
		Action: run,
		Args: []*cli.ArgDefinition{
			{
				Name:        "num1",
				Description: "first number",
				Required:    true,
			},
			{
				Name:        "num2",
				Description: "second number",
				Required:    true,
			},
		},
		UsageExamples: []*cli.UsageExample{
			{
				Command: "number-adder 1 and 2",
			},
		},
	}

	application.RunWithGlobalArgs(context.Background())
}

func run(ctx *cli.Context) error {
	number1, _ := strconv.Atoi(ctx.GetArg("number1"))
	number2, _ := strconv.Atoi(ctx.GetArg("number2"))

	fmt.Printf("%d + %d = %d", number1, number2, number1 + number2)

	return nil
}
```

## Additional features

### Generate GitHub Action configuration

Run: `go run you_main_file.go --singlecli-codegen-ga`

You also can generate in pipeline with workflow:
```yaml
name: update github action config

permissions: write-all

on:
  push:
    branches:
      - master

jobs:
  update-ga-config:
    runs-on: ubuntu-latest
    steps:
      - name: install deps
        run: sudo apt install gcc

      - name: Check out code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4 # action page: <https://github.com/actions/setup-go>
        with:
          go-version: stable

      - name: Install Go dependencies
        run: go mod download

      - name: Configure git user
        run: |
          git config user.name 'github-actions[bot]'
          git config user.email 'github-actions[bot]@users.noreply.github.com'

      - name: Hash prev version
        id: prev_version
        run: echo "hash=${{ hashFiles('action.yaml') }}" >> $GITHUB_OUTPUT

      - name: Generate new config file
        run: |
          go run ./cmd/main.go --singlecli-codegen-ga

      - name: Commit changes
        if: ${{ steps.prev_version.outputs.hash != hashFiles('action.yaml') }}
        run: |
          git add action.yaml
          git commit -m "chore: actualize ga config"
          git push
```

## Application examples
* [regexlint](https://github.com/ArtARTs36/regexlint)
* [db-exporter](https://github.com/ArtARTs36/db-exporter)
