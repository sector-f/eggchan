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
	"github.com/spf13/cobra"
)

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
	readline.PcItem("add-board"),
	readline.PcItem("list-boards"),
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
	Run: func(cmd *cobra.Command, args []string) {
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
				if runCommand(db, arguments) {
					break repl
				}
			}
		}

	},
}

func runCommand(db *sql.DB, arguments []string) (break_loop bool) {
	if len(arguments) == 0 {
		return false
	}

	var command *cobra.Command

	switch arguments[0] {
	case "add-user":
		command = addUserCommand(db)
	case "delete-user":
		command = deleteUserCommand(db)
	case "list-users":
		command = listUsersCommand(db)
	case "add-board":
		command = addBoardCommand(db)
	case "list-boards":
		command = listBoardsCommand(db)
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
