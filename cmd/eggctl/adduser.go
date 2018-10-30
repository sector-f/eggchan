package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func addUserSubcommand(ctx *cli.Context) error {
	connectionString := fmt.Sprintf("host=127.0.0.1 dbname=%s sslmode=disable", ctx.String("database"))

	var err error
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	username := ctx.Args().Get(0)
	if username != "" {
		addUserToDB(db, username)
	} else {
		return errors.New("No username provided")
	}

	return nil
}

func addUserToDB(db *sql.DB, user string) {
	passwd1, err := getPasswd("Enter password: ")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	passwd2, err := getPasswd("Enter password again: ")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	if passwd1 != passwd2 {
		fmt.Println("Passwords do not match")
		os.Exit(1)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(passwd1), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	_, err = db.Exec(
		`INSERT INTO users (username, password) VALUES ($1, $2)`,
		user,
		hashed,
	)

	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("User", user, "added successfully")
	}
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
