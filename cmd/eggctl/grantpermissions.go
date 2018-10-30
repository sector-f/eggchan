package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

func grantPermissions(ctx *cli.Context) error {
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
			`INSERT INTO user_permissions (user_id, permission) VALUES
			((SELECT id FROM users WHERE username = $1), (SELECT id FROM permissions WHERE name = $2))`,
			user,
			perm,
		)

		if err != nil {
			fmt.Printf("Error granting permission \"%s\" to %s\n", perm, user)
		}
	}

	return nil
}
