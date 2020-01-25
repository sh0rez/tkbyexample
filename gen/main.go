package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

var e = make(chan error)

func main() {
	cmdDev := &cli.Command{
		Name:  "dev",
		Usage: "Run development server",
		Action: func(c *cli.Context) error {
			if err := develop(); err != nil {
				return err
			}
			return nil
		},
	}

	cmdRender := &cli.Command{
		Name:  "render",
		Usage: "Render the markdown source for GatsbyJS",
		Action: func(c *cli.Context) error {
			return render()
		},
	}

	app := &cli.App{
		Name:        "tk by example",
		Description: "Generates https://x.tanka.dev using some custom Go code and GatsbyJS.",
		Commands:    []*cli.Command{cmdDev, cmdRender},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println()
		log.Fatalln(err)
	}
}
