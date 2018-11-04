package cmd

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func grantPermissionsCommand(db *sql.DB) *cobra.Command {
	command := cobra.Command{
		Use:           "grant-permissions",
		Short:         "Grant permissions to a user",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			permissions := args[:len(args)-1]
			username := args[len(args)-1]

			if username != "" {
				grantPermissions(db, permissions, username)
			} else {
				return errors.New("Username cannot be empty")
			}

			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: grant-permissions PERMISSIONS... USERNAME")
	})

	return &command
}

func grantPermissions(db *sql.DB, permissions []string, username string) error {
	for _, perm := range permissions {
		_, err := db.Exec(
			`INSERT INTO user_permissions (user_id, permission) VALUES
			((SELECT id FROM users WHERE username = $1), (SELECT id FROM permissions WHERE name = $2))`,
			username,
			perm,
		)

		if err != nil {
			fmt.Printf("Error granting permission \"%s\" to %s\n", perm, username)
			return err
		}
	}

	return nil
}
