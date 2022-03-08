package prompt

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/kitar0s/kyocera-ab-tool/db"
)

/*
	Driver for terminal application
*/

func Prompt(r *db.SQLiteRepository, w io.Writer) {
	l := newReadLine()
	defer l.Close()

Loop:
	for {
		p := fmt.Sprintf("%v» ", r.CurrentTable())

		l.SetPrompt(p)
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

		args := parseArgs(line)
		if  len(args) == 0 {
			continue
		}

		command := args[0]

		switch command {
		case "create_table":
			if len(args) != 2 {
				helpCommand(w, command)
				continue
			}
			tableName := args[1]
			createTable(r, w, tableName)
		case "switch_table":
			if len(args) != 2 {
				helpCommand(w, command)
				continue
			}
			tableName := args[1]
			switchTable(r, w, tableName)
		case "clear_table":
			clearTable(r, w)
		case "delete_table":
			if len(args) != 2 {
				helpCommand(w, command)
				continue
			}
			tableName := args[1]
			deleteTable(r, w, tableName)
		case "list_tables":
			listTables(r, w)
		case "show_users":
			showUsers(r, w)
		case "add_user":
			params := parseInsertArgs(line)
			if len(params) != 3 {
				helpCommand(w, command)
				continue
			}
			insertEntry(r, w, params)
		case "delete_user":
			if len(args) != 2 {
				helpCommand(w, command)
				continue
			}
			username := args[1]
			deleteEntry(r, w, username)
		case "import_csv":
			if len(args) != 2 {
				helpCommand(w, command)
				continue
			}
			f, err := os.Open(args[1])
			if err != nil {
				outputMessage(w, '-', err.Error())
				continue
			}
	
			importCSV(r, f, w)
		case "quit":
			fmt.Fprint(w, "Bye!\n\n")
			break Loop
		case "help":
			if len(args) == 2 {
				helpCommand(w, args[1])
			} else {
				listCommands(w)
			}
		default:
			msg := "type 'help' for a list of commands"
			outputMessage(w, '!', msg)
		}
	}
}

/*
	Returns customized readline instance.
*/

func newReadLine() *readline.Instance {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "",
		AutoComplete:    completions,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})
	if err != nil {
		log.Fatalf("could not create readline: %q", err)
	}

	return l
}

/*
	Returns a slice of strings that are separated by spaces
*/

func parseArgs(s string) []string {
	return strings.Fields(s)
}

/*
	Returns slice of strings from a string that is separated by commas
*/

func parseInsertArgs(s string) []string {
	sf := strings.Fields(s)
	f := strings.Join(sf[1:], " ")
	params := []string{}

	for _, str := range strings.Split(f, ",") {
		params = append(params, strings.TrimSpace(str))
	}

	return params
}

/*
	Outputs message to the user.
		+ indicates success
		- indicates error
		! indicates exclamatory
*/

func outputMessage(w io.Writer, symbol rune, msg string) {
	fmt.Fprintf(w, "[%v] %v\n\n", string(symbol), msg)
}

