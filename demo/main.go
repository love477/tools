package demo

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func CliDemo() {
	cmd1 := &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add a task to the list",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "author", Usage: "author name"},
		},
		Action: func(cCtx *cli.Context) error {
			fmt.Println("added task: ", cCtx.Args().First(), "author: ", cCtx.String("author"))
			return nil
		},
	}
	cmd2 := &cli.Command{
		Name:    "complete",
		Aliases: []string{"c"},
		Usage:   "complete a task on the list",
		Action: func(cCtx *cli.Context) error {
			fmt.Println("completed task: ", cCtx.Args().First())
			return nil
		},
	}

	cmd3 := &cli.Command{
		Name:        "template",
		Aliases:     []string{"t"},
		Usage:       "options for task templates",
		Subcommands: []*cli.Command{cmd1, cmd2},
	}

	app := &cli.App{
		Name:     "demo",
		Usage:    "it is demo for cli",
		Commands: []*cli.Command{cmd1, cmd2, cmd3},
		Action: func(cCtx *cli.Context) error {
			fmt.Println("completed task: ", cCtx.Args())
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
