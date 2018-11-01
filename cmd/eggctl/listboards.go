package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

func listBoardsCommand() cli.Command {
	return cli.Command{
		Name:  "list-boards",
		Usage: "List boards in the database",
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
			return listBoards(ctx)
		},
	}
}

// TODO: add bump_limit, post_limit, max_num_threads
type board struct {
	name        string
	description sql.NullString
}

func listBoards(ctx *cli.Context) error {
	connectionString := fmt.Sprintf("host=127.0.0.1 dbname=%s sslmode=disable", ctx.String("database"))

	var err error
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	rows, err := db.Query(`SELECT name, description FROM boards ORDER BY id ASC`)
	if err != nil {
		return err
	}

	boardList := []board{}
	for i := 0; rows.Next(); i++ {
		var b board
		if err := rows.Scan(&b.name, &b.description); err != nil {
			return err
		}
		boardList = append(boardList, b)
	}

	for _, board := range boardList {
		if board.description.Valid {
			fmt.Printf("%s - %s\n", board.name, board.description.String)
		} else {
			fmt.Printf("%s\n", board.name)
		}
	}

	return nil
}
