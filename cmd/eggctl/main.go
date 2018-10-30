package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:  "add-user",
			Usage: "Add a new user to the database",
			Flags: []cli.Flag{cli.StringFlag{
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
			},
			Action: func(ctx *cli.Context) error {
				return addUser(ctx)
			},
		},
		{
			Name:  "list-permissions",
			Usage: "List available permissions",
			Flags: []cli.Flag{cli.StringFlag{
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
			},
			Action: func(ctx *cli.Context) error {
				return listPermissions(ctx)
			},
		},
		{
			Name:  "grant-permissions",
			Usage: "Grant permissions to a user",
			Flags: []cli.Flag{cli.StringFlag{
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
			},
			Action: func(ctx *cli.Context) error {
				return grantPermissions(ctx)
			},
		},
		{
			Name:  "revoke-permissions",
			Usage: "Revoke permissions from a user",
			Flags: []cli.Flag{cli.StringFlag{
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
			},
			Action: func(ctx *cli.Context) error {
				return revokePermissions(ctx)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
