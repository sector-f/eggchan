package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func addUserCommand(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:           "add-user",
		Short:         "Add user to the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			if username != "" {
				if err := addUserToDB(db, username); err != nil {
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
}

func addUserToDB(db *sql.DB, user string) error {
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

	hashed, err := bcrypt.GenerateFromPassword([]byte(passwd1), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`INSERT INTO users (username, password) VALUES ($1, $2)`,
		user,
		hashed,
	)

	if err != nil {
		return err
	}

	return nil
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
