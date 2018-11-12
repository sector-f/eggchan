package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sector-f/eggchan"
)

func grantPermissionsCommand() *cobra.Command {
	command := cobra.Command{
		Use:           "grant-permissions",
		Short:         "Grant permissions to a user",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			perms_string := args[:len(args)-1]
			username := args[len(args)-1]

			permissions := []eggchan.Permission{}
			for _, perm := range perms_string {
				permissions = append(permissions, eggchan.Permission{perm})
			}

			if username != "" {
				Service.GrantPermissions(username, permissions)
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
