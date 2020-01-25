package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	cli "github.com/urfave/cli/v2"
)

var e = make(chan error)

func main() {
	app := &cli.App{
		Name:        "tk by example",
		Description: "Generates https://x.tanka.dev using some custom Go code and GatsbyJS.",
	}

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

	cmdCommands := &cli.Command{
		Name:   "commands",
		Hidden: true,
		Action: func(c *cli.Context) error {
			s := ""
			for _, c := range app.Commands {
				if c.Hidden {
					continue
				}
				s += c.Name + " "
			}
			fmt.Println(strings.TrimSuffix(s, " "))
			return nil
		},
	}

	app.Commands = []*cli.Command{cmdDev, cmdRender, cmdCommands}

	if err := app.Run(os.Args); err != nil {
		fmt.Println()
		log.Fatalln(err)
	}
}
