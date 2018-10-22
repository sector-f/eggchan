package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "eggchan"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bind, b",
			Value: "127.0.0.1:8000",
			Usage: "Address/port to bind to",
		},
		cli.StringFlag{
			Name:   "database, d",
			Usage:  "Database name",
			EnvVar: "EGGCHAN_DB_NAME",
		},
		cli.StringFlag{
			Name:   "username, u",
			Usage:  "Database username",
			EnvVar: "EGGCHAN_DB_USERNAME",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "Database password",
			EnvVar: "EGGCHAN_DB_PASSWORD",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		server := Server{}
		server.Initialize(
			ctx.String("username"),
			ctx.String("password"),
			ctx.String("database"),
		)
		server.Run(ctx.String("bind"))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
