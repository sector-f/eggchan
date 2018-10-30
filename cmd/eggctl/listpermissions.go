package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

func listPermissionsCommand() cli.Command {
	return cli.Command{
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
	}
}

func listPermissions(ctx *cli.Context) error {
	connectionString := fmt.Sprintf("host=127.0.0.1 dbname=%s sslmode=disable", ctx.String("database"))

	var err error
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	rows, err := db.Query(`SELECT name FROM permissions`)
	if err != nil {
		return err
	}

	perms := []string{}
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return err
		}
		perms = append(perms, p)
	}

	fmt.Println(strings.Join(perms, "\n"))

	return nil
}
