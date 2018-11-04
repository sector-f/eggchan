package cmd

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func revokePermissionsCommand(db *sql.DB) *cobra.Command {
	command := cobra.Command{
		Use:           "revoke-permissions",
		Short:         "Revoke permissions from a user",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			permissions := args[:len(args)-1]
			username := args[len(args)-1]

			if username != "" {
				if err := revokePermissions(db, permissions, username); err != nil {
					return err
				}
			} else {
				return errors.New("Username cannot be empty")
			}

			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: revoke-permissions PERMISSIONS... USERNAME")
	})

	return &command
}

func revokePermissions(db *sql.DB, permissions []string, username string) error {
	for _, perm := range permissions {
		result, err := db.Exec(
			`DELETE FROM user_permissions
			WHERE user_id = (SELECT id FROM users WHERE username = $1)
			AND permission = (SELECT id FROM permissions WHERE name = $2)`,
			username,
			perm,
		)

		if err != nil {
			fmt.Printf("Error revoking permission \"%s\" from %s\n", perm, username)
		}

		affected, _ := result.RowsAffected()
		if affected != 1 {
			fmt.Printf("Error revoking permission \"%s\" from %s\n", perm, username)
		}
	}

	return nil
}
