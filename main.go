package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "composer-repo"
	app.Usage = "sync multiple git repos"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "initialize repo with manifest.json and docker-compose.yml",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "repo,r",
					Usage: "path to manifest repo",
				},
			},
			Action: func(c *cli.Context) {},
		},
		{
			Name:   "sync",
			Usage:  "syncronize repos from manifest.json",
			Action: func(c *cli.Context) {},
		},
	}

	app.Run(os.Args)
}
