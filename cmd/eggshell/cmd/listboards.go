package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listBoardsCommand() *cobra.Command {
	command := cobra.Command{
		Use:           "list-boards",
		Short:         "List the boards in the Eggchan database",
		SilenceErrors: false,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			boards, err := Service.ListBoards()
			if err != nil {
				return err
			}

			for _, board := range boards {
				fmt.Printf("%s - %s\n", board.Name, board.Description)
			}

			return nil
		},
	}

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: list-boards")
	})

	return &command
}
