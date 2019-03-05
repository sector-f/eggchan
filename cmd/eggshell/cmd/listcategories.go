package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listCategoriesCommand() *cobra.Command {
	command := cobra.Command{
		Use:           "list-categories",
		Short:         "List the categories in the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			categories, err := Service.ListCategories()
			if err != nil {
				return err
			}

			for _, category := range categories {
				fmt.Printf("%s\n", category.Name)
			}

			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: list-categories")
	})

	return &command
}
