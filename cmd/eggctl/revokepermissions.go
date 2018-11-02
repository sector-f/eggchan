package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

func revokePermissionsCommand() cli.Command {
	return cli.Command{
		Name:  "revoke-permissions",
		Usage: "Revoke permissions from a user",
		Flags: []cli.Flag{
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
		},
		Action: func(ctx *cli.Context) error {
			return revokePermissions(ctx)
		},
	}
}

func revokePermissions(ctx *cli.Context) error {
	connectionString := fmt.Sprintf("host=127.0.0.1 dbname=%s sslmode=disable", ctx.String("database"))

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	num_args := len(ctx.Args())
	if num_args < 2 {
		return errors.New("Not enough arguments provided")
	}
	user := ctx.Args()[num_args-1]

	for _, perm := range ctx.Args()[:num_args-1] {
		_, err := db.Exec(
			`DELETE FROM user_permissions
			WHERE user_id = (SELECT id FROM users WHERE username = $1)
			AND permission = (SELECT id FROM permissions WHERE name = $2)`,
			user,
			perm,
		)
		if err != nil {
			fmt.Printf("Error revoking permission \"%s\" from %s\n", perm, user)
		}
	}

	return nil
}
