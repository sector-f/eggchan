package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var Description string
var Category string

func addBoardCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "add-board",
		Short:         "Add board to the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			board := args[0]

			description, _ := cmd.Flags().GetString("description")
			category, _ := cmd.Flags().GetString("category")

			if board != "" {
				if err := Service.AddBoard(board, description, category); err != nil {
					fmt.Printf("Error: %s\n", err)
				} else {
					fmt.Println("Board", board, "added successfully")
				}
			} else {
				return errors.New("Board name cannot be empty")
			}

			return nil
		},
	}

	command.Flags().StringVarP(&Description, "description", "d", "", "Board description")
	command.Flags().StringVarP(&Category, "category", "c", "", "Board category")
	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: add-board [--description DESCRIPTION] [--category CATEGORY] BOARDNAME")
	})

	return command
}
