package cmd

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func listPermissionsCommand(db *sql.DB) *cobra.Command {
	command := cobra.Command{
		Use:           "list-permissions",
		Short:         "List available permissions",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := listPermissions(db); err != nil {
				return err
			}

			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: list-permissions")
	})

	return &command
}

func listPermissions(db *sql.DB) error {
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

	fmt.Println(strings.Join(perms, "\n"))

	return nil
}
