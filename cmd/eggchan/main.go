package main

import (
	"log"
	"os"

	"github.com/sector-f/eggchan/database/postgres"
	"github.com/sector-f/eggchan/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "eggchan"
	app.Version = "0.1.0"
	app.Usage = "A headless JSON textboard"

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
		pgOptions := postgres.Options{Hostname: "127.0.0.1", Database: ctx.String("database")}
		service, err := postgres.New(pgOptions)
		if err != nil {
			return err
		}

		httpServer := server.HttpServer{
			BoardService: service,
			AdminService: service,
			AuthService:  service,
		}
		httpServer.Initialize()
		httpServer.Run(ctx.String("bind"))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
