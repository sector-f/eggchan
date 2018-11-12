package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func deleteUserCommand() *cobra.Command {
	command := cobra.Command{
		Use:           "delete-user",
		Short:         "Delete user from the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			if username != "" {
				if err := Service.DeleteUser(username); err != nil {
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
