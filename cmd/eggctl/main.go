package main

import (
	"fmt"
	"os"
	"syscall"

	// "github.com/urfave/cli"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
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

	fmt.Println(string(hashed))
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
