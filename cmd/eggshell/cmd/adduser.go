package cmd

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func addUserCommand() *cobra.Command {
	command := cobra.Command{
		Use:           "add-user",
		Short:         "Add user to the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			if username != "" {
				passwd1, err := getPasswd("Enter password: ")
				if err != nil {
					return err
				}

				passwd2, err := getPasswd("Enter password again: ")
				if err != nil {
					return err
				}

				if passwd1 != passwd2 {
					return errors.New("Passwords do not match")
				}

				if err := Service.AddUser(username, passwd1); err != nil {
					fmt.Printf("Error: %s\n", err)
				} else {
					fmt.Println("User", username, "added successfully")
				}
			} else {
				return errors.New("Username cannot be empty")
			}

			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: add-user USERNAME")
	})

	return &command
}

func getPasswd(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}

	return string(bytePassword), nil
}
