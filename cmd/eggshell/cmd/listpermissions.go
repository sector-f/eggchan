package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func listPermissionsCommand() *cobra.Command {
	command := cobra.Command{
		Use:           "list-permissions",
		Short:         "List available permissions",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			perms, err := Service.ListPermissions()
			if err != nil {
				return err
			}

			perm_names := []string{}
			for _, perm := range perms {
				perm_names = append(perm_names, perm.Name)
			}

			fmt.Println(strings.Join(perm_names, "\n"))

			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: list-permissions")
	})

	return &command
}
