package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

func listUsersCommand() cli.Command {
	return cli.Command{
		Name:  "list-users",
		Usage: "List users in the database",
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
			return listUsers(ctx)
		},
	}
}

type user struct {
	name  string
	perms []string
}

func listUsers(ctx *cli.Context) error {
	connectionString := fmt.Sprintf("host=127.0.0.1 dbname=%s sslmode=disable", ctx.String("database"))
	userList := make(map[string][]string)

	var err error
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	rows, err := db.Query(`SELECT username FROM users ORDER BY id ASC`)
	if err != nil {
		return err
	}

	users := []string{}
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return err
		}
		users = append(users, u)
	}

	for _, user := range users {
		rows, err = db.Query(
			`SELECT name FROM permissions p
			INNER JOIN user_permissions up ON p.id = up.permission
			INNER JOIN users u ON u.id = up.user_id
			WHERE u.username = $1
			ORDER BY p.id ASC`,
			user,
		)
		if err != nil {
			return err
		}

		permissions := []string{}
		for rows.Next() {
			var p string
			if err := rows.Scan(&p); err != nil {
				return err
			}
			permissions = append(permissions, p)
		}

		userList[user] = permissions
	}

	for k, v := range userList {
		fmt.Printf("%s: %s\n", k, strings.Join(v, " "))
	}

	return nil
}
