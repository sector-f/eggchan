package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		addUserCommand(),
		listUsersCommand(),
		deleteUserCommand(),
		addBoardCommand(),
		listBoardsCommand(),
		listPermissionsCommand(),
		grantPermissionsCommand(),
		revokePermissionsCommand(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
