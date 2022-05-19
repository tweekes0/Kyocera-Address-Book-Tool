package prompt

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/tweekes0/kyocera-ab-tool/db"
)

/*
	Driver for terminal application
*/

func Prompt(r *db.SQLiteRepository, rd io.ReadCloser, w io.Writer) {
	l := newReadLine(rd)
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
			case "export_table":
				all, err := r.All()
				if err != nil {
					OutputMessage(w, '-', err.Error())
				}

				if len(all) == 0 {
					msg := "cannot export empty table"
					OutputMessage(w, '-', msg)
					continue
				}

				f, err := createFile(r.CurrentTable())
				if err != nil {
					OutputMessage(w, '-', err.Error())
					continue
				}

				exportTable(r, w, f)
				f.Close()
			case "help":
				listCommands(w)
			case "exit", "quit":
				break Loop
			case "create_table", "switch_table", "delete_table", "add_user",
				"delete_user", "update_user", "import_csv":
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
			case "update_user":
				updateUser(r, w, param)
			case "import_csv":
				f, err := os.Open(param)
				if err != nil {
					OutputMessage(w, '-', err.Error())
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
