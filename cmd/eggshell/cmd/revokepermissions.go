package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sector-f/eggchan"
)

func revokePermissionsCommand() *cobra.Command {
	command := cobra.Command{
		Use:           "revoke-permissions",
		Short:         "Revoke permissions from a user",
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
				if err := Service.RevokePermissions(username, permissions); err != nil {
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
