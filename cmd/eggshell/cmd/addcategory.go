package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func addCategoryCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "add-category",
		Short:         "Add category to the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			category := args[0]

			if category != "" {
				if err := Service.AddCategory(category); err != nil {
					fmt.Printf("Error: %s\n", err)
				} else {
					fmt.Println("Category", category, "added successfully")
				}
			} else {
				return errors.New("Category name cannot be empty")
			}

			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: add-category CATEGORYNAME")
	})

	return command
}
