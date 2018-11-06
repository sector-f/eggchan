package cmd

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func listBoardsCommand(db *sql.DB) *cobra.Command {
	command := cobra.Command{
		Use:           "list-boards",
		Short:         "List the boards in the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := listBoards(db)
			if err != nil {
				return err
			}
			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: list-boards")
	})

	return &command
}

// TODO: add bump_limit, post_limit, max_num_threads
type board struct {
	name        string
	description sql.NullString
}

func listBoards(db *sql.DB) error {
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
