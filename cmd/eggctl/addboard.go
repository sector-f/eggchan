package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

func addBoardCommand() cli.Command {
	return cli.Command{
		Name:  "add-board",
		Usage: "Add a new board to the database",
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
			cli.StringFlag{
				Name:  "description",
				Usage: "Board description",
			},
			cli.StringFlag{
				Name:  "category",
				Usage: "Board category",
			},
		},
		Action: func(ctx *cli.Context) error {
			return addBoard(ctx)
		},
	}
}

func addBoard(ctx *cli.Context) error {
	connectionString := fmt.Sprintf("host=127.0.0.1 dbname=%s sslmode=disable", ctx.String("database"))

	var err error
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	boardname := strings.TrimSpace(ctx.Args().Get(0))
	description := strings.TrimSpace(ctx.String("description"))
	category := strings.TrimSpace(ctx.String("category"))
	if boardname != "" {
		addBoardToDB(db, boardname, description, category)
	} else {
		return errors.New("No board name provided")
	}

	return nil
}

func addBoardToDB(db *sql.DB, name string, description string, category string) error {
	var d sql.NullString
	if description == "" {
		d = sql.NullString{"", false}
	} else {
		d = sql.NullString{description, true}
	}

	var c sql.NullString
	if category == "" {
		c = sql.NullString{"", false}
	} else {
		c = sql.NullString{category, true}
	}

	_, err := db.Exec(
		`INSERT INTO boards (name, description, category) VALUES ($1, $2, (SELECT id FROM categories WHERE name = $3))`,
		name,
		d,
		c,
	)

	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Board added successfully")
	}
	return nil
}
