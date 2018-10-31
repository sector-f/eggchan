package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

func deleteUserCommand() cli.Command {
	return cli.Command{
		Name:  "delete-user",
		Usage: "Remove a user from the database",
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
			return deleteUser(ctx)
		},
	}
}

func deleteUser(ctx *cli.Context) error {
	connectionString := fmt.Sprintf("host=127.0.0.1 dbname=%s sslmode=disable", ctx.String("database"))

	var err error
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	username := ctx.Args().Get(0)
	if username != "" {
		deleteUserFromDB(db, username)
	} else {
		return errors.New("No username provided")
	}

	return nil
}

func deleteUserFromDB(db *sql.DB, user string) {
	_, err := db.Exec(
		`DELETE FROM users WHERE username = $1`,
		user,
	)

	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("User", user, "deleted successfully")
	}
}
