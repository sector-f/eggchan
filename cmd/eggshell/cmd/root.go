package cmd

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	shellquote "github.com/kballard/go-shellquote"
	_ "github.com/lib/pq"
	"github.com/sector-f/eggchan/database/postgres"
	"github.com/spf13/cobra"
)

var Service postgres.EggchanService

// Cobra global variables
var Database string
var Username string
var Password string
var Egg bool

// Readline tab completion
var completer = readline.NewPrefixCompleter(
	readline.PcItem("add-user"),
	readline.PcItem("delete-user"),
	readline.PcItem("list-users"),
	readline.PcItem("grant-permissions"),
	readline.PcItem("revoke-permissions"),
	readline.PcItem("list-permissions"),
	readline.PcItem("add-board"),
	readline.PcItem("list-boards"),
	readline.PcItem("add-category"),
	readline.PcItem("list-categories"),
	readline.PcItem("help"),
	readline.PcItem("exit"),
)

func init() {
	rootCmd.PersistentFlags().StringVar(&Database, "database", "eggchan", "Database name")
	rootCmd.PersistentFlags().StringVar(&Username, "username", "eggchan", "Database username")
	rootCmd.PersistentFlags().StringVar(&Password, "password", "", "Database password")
	rootCmd.PersistentFlags().BoolVar(&Egg, "egg", false, "Enable egg")

	rootCmd.PersistentFlags().MarkHidden("egg")
}

var rootCmd = &cobra.Command{
	Use:   "eggshell",
	Short: "Command-line interface to the Eggchan database",
	PreRun: func(cmd *cobra.Command, args []string) {
		connectionString := fmt.Sprintf("host=127.0.0.1 dbname=%s sslmode=disable", cmd.Flag("database").Value.String())
		var err error
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			fmt.Printf("Error establishing database connection: %s\n", err)
			return
		}

		err = db.Ping()
		if err != nil {
			fmt.Printf("Error establishing database connection: %s\n", err)
			return
		}

		Service = postgres.EggchanService{db}

		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		var prompt string
		if egg, _ := cmd.Flags().GetBool("egg"); egg {
			prompt = "🥚 "
		} else {
			prompt = "> "
		}

		l, err := readline.NewEx(&readline.Config{
			Prompt:          prompt,
			AutoComplete:    completer,
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
		})
		if err != nil {
			panic(err)
		}
		defer l.Close()

	repl:
		for {
			line, err := l.Readline()
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}

			arguments, err := shellquote.Split(strings.TrimSpace(line))
			if err != nil {
				fmt.Printf("Syntax error: %s\n", err)
				continue
			} else {
				if runCommand(arguments) {
					break repl
				}
			}
		}

	},
}

func runCommand(arguments []string) (break_loop bool) {
	if len(arguments) == 0 {
		return false
	}

	var command *cobra.Command

	switch arguments[0] {
	case "add-user":
		command = addUserCommand()
	case "delete-user":
		command = deleteUserCommand()
	case "list-users":
		command = listUsersCommand()
	case "grant-permissions":
		command = grantPermissionsCommand()
	case "revoke-permissions":
		command = revokePermissionsCommand()
	case "list-permissions":
		command = listPermissionsCommand()
	case "add-board":
		command = addBoardCommand()
	case "list-boards":
		command = listBoardsCommand()
	case "add-category":
		command = addCategoryCommand()
	case "list-categories":
		command = listCategoriesCommand()
	case "help":
		commands := []string{
			"add-user",
			"delete-user",
			"list-users",
			"grant-permissions",
			"revoke-permissions",
			"list-permissions",
			"add-board",
			"list-boards",
			"add-category",
			"list-categories",
			"help",
			"exit",
		}
		fmt.Println("Available commands:")
		fmt.Println(strings.Join(commands, "\n"))
		return false
	case "exit":
		return true
	default:
		fmt.Printf("Error: Unknown command \"%s\"\n", arguments[0])
		return false
	}

	command.SetArgs(arguments[1:])
	command.Execute()

	return false
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
