package cmd

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var Description string
var Category string

func addBoardCommand(db *sql.DB) *cobra.Command {
	command := &cobra.Command{
		Use:           "add-board",
		Short:         "Add board to the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			board := args[0]

			description, _ := cmd.Flags().GetString("description")
			category, _ := cmd.Flags().GetString("category")

			if board != "" {
				if err := addBoardToDB(db, board, description, category); err != nil {
					fmt.Printf("Error: %s\n", err)
				} else {
					fmt.Println("Board", board, "added successfully")
				}
			} else {
				return errors.New("Board name cannot be empty")
			}

			return nil
		},
	}

	command.Flags().StringVarP(&Description, "description", "d", "", "Board description")
	command.Flags().StringVarP(&Category, "category", "c", "", "Board category")
	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: add-board [--description DESCRIPTION] [--category CATEGORY] BOARDNAME")
	})

	return command
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
		return err
	}

	return nil
}
