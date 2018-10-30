package main

import (
	"database/sql"
	// "errors"
	"fmt"
	"strings"
	// "os"
	// "syscall"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

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

	fmt.Println("Available permissions:", strings.Join(perms, ", "))

	return nil
}
