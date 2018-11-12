package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func listUsersCommand() *cobra.Command {
	command := cobra.Command{
		Use:           "list-users",
		Short:         "List the users in the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			users, err := Service.ListUsers()
			if err != nil {
				return err
			}

			for _, user := range users {
				if len(user.Perms) == 0 {
					fmt.Printf("%s:\n", user.Name)
				} else {
					fmt.Printf("%s: %s\n", user.Name, strings.Join(user.Perms, ", "))
				}
			}

			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: list-users")
	})

	return &command
}
