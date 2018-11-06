package cmd

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func deleteUserCommand(db *sql.DB) *cobra.Command {
	command := cobra.Command{
		Use:           "delete-user",
		Short:         "Delete user from the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			if username != "" {
				if err := deleteUser(db, username); err != nil {
					fmt.Printf("Error: %s\n", err)
				} else {
					fmt.Println("User", username, "deleted successfully")
				}
			} else {
				return errors.New("Username cannot be empty")
			}

			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: delete-user USERNAME")
	})

	return &command
}

func deleteUser(db *sql.DB, user string) error {
	result, err := db.Exec(
		`DELETE FROM users WHERE username = $1`,
		user,
	)

	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected != 1 {
		return errors.New("User not found")
	}

	return nil
}
