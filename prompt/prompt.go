package prompt

import (
	// "bufio"
	"fmt"
	"io"
	"log"

	"strings"

	"github.com/chzyer/readline"
	"github.com/kitar0s/kyocera-ab-tool/db"
)

/*

 */
func Prompt(r *db.SQLiteRepository, rd io.Reader, w io.Writer) {
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
		command := args[0]

		switch command {
		case "create_table":
			if len(args) != 2 {
				continue
			}
			tableName := args[1]
			createTable(r, w, tableName)
		case "switch_table":
			if len(args) != 2 {
				continue
			}
			tableName := args[1]
			switchTable(r, w, tableName)
		case "users":
			showUsers(r, w)
		case "add", "insert":
			params := parseInsertArgs(line)
			if len(params) != 3 {
				continue
			}
			insertEntry(r, w, params)
		case "delete_user":
			if len(args) != 2 {
				continue
			}
			username := args[1]
			deleteEntry(r, w, username)
		case "q", "quit": 
			fmt.Fprint(w, "Bye!\n\n")
			break Loop
		default:
			fmt.Fprint(w, "[-] enter 'help' or 'h' for a list of commands\n\n")
		}
	}
}

func newReadLine() *readline.Instance {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31mprompt»\033[0m ",
		AutoComplete:    completions,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
	})
	if err != nil {
		log.Fatalf("could not create readline: %q", err)
	}

	return l
}

func parseArgs(s string) []string {
	return strings.Fields(s)
}

func parseInsertArgs(s string) []string {
	sf := strings.Fields(s)
	f := strings.Join(sf[1:], " ")
	params := []string{}

	for _, str := range strings.Split(f, ",") {
		params = append(params, strings.TrimSpace(str))	
	}
	
	return params
}

