package cmd

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func listUsersCommand(db *sql.DB) *cobra.Command {
	command := cobra.Command{
		Use:           "list-users",
		Short:         "List the users in the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := listUsers(db)
			if err != nil {
				return err
			}
			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: list-users")
	})

	return &command
}

type user struct {
	name  string
	perms []string
}

func listUsers(db *sql.DB) error {
	userList := [][]string{}

	rows, err := db.Query(`SELECT username FROM users ORDER BY id ASC`)
	if err != nil {
		return err
	}

	for i := 0; rows.Next(); i++ {
		var u string
		if err := rows.Scan(&u); err != nil {
			return err
		}
		userList = append(userList, []string{u})
	}

	for i, user := range userList {
		rows, err = db.Query(
			`SELECT name FROM permissions p
			INNER JOIN user_permissions up ON p.id = up.permission
			INNER JOIN users u ON u.id = up.user_id
			WHERE u.username = $1
			ORDER BY p.id ASC`,
			user[0],
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

		userList[i] = append(userList[i], permissions...)
	}

	for _, user := range userList {
		if len(user) > 1 {
			fmt.Printf("%s: %s\n", user[0], strings.Join(user[1:], " "))
		} else {
			fmt.Printf("%s:\n", user[0])
		}
	}

	return nil
}
