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

		command, param := parseArgs(line)
		if command == "" {
			continue
		}

		switch {
		case param == "":
			switch command {
			case "clear_table":
				clearTable(r, w)
			case "list_tables":
				listTables(r, w)
			case "show_users":
				showUsers(r, w)
			case "help":
				listCommands(w)
			case "quit":
				break Loop
			case "create_table", "switch_table", "delete_table", "add_user",
				"delete_user", "import_csv":
				helpCommand(w, command)
			default:
				helpUser(w)

			}

		case param != "":
			switch command {
			case "create_table":
				createTable(r, w, param)
			case "switch_table":
				switchTable(r, w, param)
			case "delete_table":
				deleteTable(r, w, param)
			case "add_user":
				addUser(r, w, param)
			case "delete_user":
				deleteUser(r, w, param)
			case "import_csv":
				f, err := os.Open(param)
				if err != nil {
					outputMessage(w, '-', err.Error())
					continue
				}
				importCSV(r, f, w)
			case "help":
				helpCommand(w, param)
			default:
				helpUser(w)
			}
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
	Outputs message to the user.
		+ indicates success
		- indicates error
		! indicates exclamatory
*/

func outputMessage(w io.Writer, symbol rune, msg string) {
	fmt.Fprintf(w, "[%v] %v\n\n", string(symbol), msg)
}

/*
	Remove the quotes from the supplied string
*/

func stripQuotes(s string) string {
	s = strings.Replace(s, `"`, "", -1)
	s = strings.Replace(s, `'`, "", -1)
	s = strings.Replace(s, "`", "", -1)

	return s
}

/*
	Take a from user input and returns two strings a command and it's optional
	parameter
*/

func parseArgs(s string) (command string, param string) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		command = ""
		param = ""
		return
	}

	if len(fields) == 1 {
		command = fields[0]
		param = ""
		return
	}

	command = fields[0]
	param = strings.Join(fields[1:], " ")
	param = stripQuotes(param)

	return
}

func helpUser(w io.Writer) {
	msg := "type 'help' for a list of commands"
	outputMessage(w, '!', msg)
}
