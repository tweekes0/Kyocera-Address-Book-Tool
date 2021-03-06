package prompt

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/chzyer/readline"
	"github.com/rodaine/table"
	"github.com/tweekes0/kyocera-ab-tool/db"
	"github.com/tweekes0/kyocera-ab-tool/exporter"
	"github.com/tweekes0/kyocera-ab-tool/importer"
)

/*
	List of PcItems that readline structs use for terminal autocompletes.
*/

var completions = readline.NewPrefixCompleter(
	readline.PcItem("create_table"),
	readline.PcItem("switch_table"),
	readline.PcItem("clear_table"),
	readline.PcItem("delete_table"),
	readline.PcItem("export_table"),
	readline.PcItem("list_tables"),
	readline.PcItem("show_users"),
	readline.PcItem("add_user"),
	readline.PcItem("delete_user"),
	readline.PcItem("update_user"),
	readline.PcItem("import_csv"),
	readline.PcItem("exit"),

	readline.PcItem("help",
		readline.PcItem("create_table"),
		readline.PcItem("switch_table"),
		readline.PcItem("clear_table"),
		readline.PcItem("delete_table"),
		readline.PcItem("export_table"),
		readline.PcItem("list_tables"),
		readline.PcItem("show_users"),
		readline.PcItem("add_user"),
		readline.PcItem("delete_user"),
		readline.PcItem("update_user"),
		readline.PcItem("import_csv"),
		readline.PcItem("exit"),
	),
)

/*
	Map of commands with descriptions and usage for printing help information
*/

var commands = map[string]struct {
	description string
	usage       string
}{
	"create_table": {
		description: "creates new table and sets it to the current table",
		usage:       "create_table 'TABLE_NAME'",
	},
	"switch_table": {
		description: "switch the current table",
		usage:       "switch_table 'TABLE_NAME'",
	},
	"clear_table": {
		description: "clears all users from the current table",
		usage:       "clear_table",
	},
	"delete_table": {
		description: "deletes the specified table",
		usage:       "delete_table 'TABLE_NAME'",
	},
	"export_table": {
		description: "exports the current table to an xml file in the Address Books directory",
		usage:       "export_table",
	},
	"list_tables": {
		description: "list all tables",
		usage:       "list_tables",
	},
	"show_users": {
		description: "show all the users in the current table",
		usage:       "show_users",
	},
	"add_user": {
		description: "add user to the current table. Fields must be separated by commas",
		usage:       "add_user 'NAME,USERNAME,EMAIL'",
	},
	"delete_user": {
		description: "delete a single user from the current table",
		usage:       "delete_user 'USERNAME'",
	},
	"update_user": {
		description: "update user in the current table. Fields must be separated by commas",
		usage:       "update_user 'USERNAME' 'NAME,USERNAME,EMAIL'",
	},
	"import_csv": {
		description: "import users from csv file into current table",
		usage:       "import_csv 'PATH_TO_FILE'",
	},
	"exit": {
		description: "exits the program",
		usage:       "exit",
	},
}

/*
	Function that will print the description and usage of a command if it exists,
	otherwise it will list the available commands.
*/

func helpCommand(w io.Writer, s string) {
	if command, ok := commands[s]; ok {
		fmt.Fprintf(w, "\n%v", command.description)
		fmt.Fprintf(w, "\nusage: %v\n\n", command.usage)
	} else {
		listCommands(w)
	}
}

/*
	List all the commands to the user.
*/

func listCommands(w io.Writer) {
	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	fmt.Fprint(w, "\nCommands:\n")
	for _, k := range keys {
		fmt.Fprintf(w, "     %-15v : %10v\n", k, commands[k].description)
	}

	fmt.Fprint(w, "\n")
}

/*
	Create a new table and write to io.Writer the success or failure
	of the operation.
*/

func createTable(r *db.SQLiteRepository, w io.Writer, tableName string) {
	err := r.NewTable(tableName)
	if err != nil {
		OutputMessage(w, '-', err.Error())
	} else {
		msg := fmt.Sprintf("%v was created successfully", r.CurrentTable())
		OutputMessage(w, '+', msg)
	}
}

/*
	Switch to a table and write to w if the operation fails.
*/

func switchTable(r *db.SQLiteRepository, w io.Writer, tableName string) {
	err := r.SwitchTable(tableName)
	if err != nil {
		OutputMessage(w, '-', err.Error())
	}
}

/*
	Display all the users that in the current table.
*/

func showUsers(r *db.SQLiteRepository, w io.Writer) {
	all, err := r.All()
	switch {
	case err != nil:
		OutputMessage(w, '-', err.Error())
	case len(all) == 0:
		msg := fmt.Sprintf("%v is empty", r.CurrentTable())
		OutputMessage(w, '!', msg)
	default:
		msg := fmt.Sprintf("users of %v", r.CurrentTable())
		OutputMessage(w, '+', msg)

		all, err := r.All()
		if err != nil {
			OutputMessage(w, '-', err.Error())
			return
		}

		tbl := table.New("ID", "Name", "Username", "Email")		

		for _, entry := range all {
			tbl.AddRow(entry.ID, entry.Name, entry.Username, entry.Email)
		}

		tbl.Print()
	}
}

/*
	Inserts a new user's Entry into the current table, granted that the
	params are valid
*/

func addUser(r *db.SQLiteRepository, w io.Writer, params string) {
	fields := strings.Split(params, ",")
	if len(fields) != 3 {
		msg := "invalid number of fields"
		OutputMessage(w, '-', msg)
		return
	}

	for i, field := range fields {
		fields[i] = strings.TrimSpace(field)
	}

	e, err := db.NewEntry(fields[0], fields[1], fields[2])
	if err != nil {
		OutputMessage(w, '-', err.Error())
		return
	} else {
		_, err = r.Insert(*e)
		if err != nil {
			OutputMessage(w, '-', err.Error())
			return
		} else {
			msg := fmt.Sprintf("%v was added successfully", e.Name)
			OutputMessage(w, '+', msg)
		}
	}
}

func updateUser(r *db.SQLiteRepository, w io.Writer, params string) {
	p := strings.Split(params, " ")
	_, err := r.GetByUsername(p[0])

	if err != nil {
		OutputMessage(w, '-', err.Error())
		return
	}

	fields := strings.Split(strings.Join(p[1:], " "), ",")
	if len(fields) != 3 {
		msg := "invalid number of fields"
		OutputMessage(w, '-', msg)
		return
	}

	for i, field := range fields {
		fields[i] = strings.TrimSpace(field)
	}

	e, err := db.NewEntry(fields[0], fields[1], fields[2])
	if err != nil {
		OutputMessage(w, '-', err.Error())
		return
	}

	_, err = r.Update(p[0], e)
	if err != nil {
		OutputMessage(w, '-', err.Error())
		return
	}

	msg := fmt.Sprintf("%v has been updated", e.Name)
	OutputMessage(w, '+', msg)
}

/*
	Delete an user's Entry from the database given a valid username.
*/

func deleteUser(r *db.SQLiteRepository, w io.Writer, username string) {
	e, err := r.GetByUsername(username)
	if err != nil {
		OutputMessage(w, '-', err.Error())
		return
	}

	err = r.Delete(username)
	if err != nil {
		OutputMessage(w, '-', err.Error())
	} else {
		msg := fmt.Sprintf("%v was deleted successfully", e.Name)
		OutputMessage(w, '+', msg)
	}
}

/*
	Clear all the entries from the current table
*/

func clearTable(r *db.SQLiteRepository, w io.Writer) {
	r.ClearTable()

	msg := fmt.Sprintf("%v was cleared successfully", r.CurrentTable())
	OutputMessage(w, '+', msg)
}

/*
	Deletes the specified table. DEFAULT_TABLE cannot be deleted.
*/

func deleteTable(r *db.SQLiteRepository, w io.Writer, tableName string) {
	err := r.DeleteTable(tableName)
	if err != nil {
		OutputMessage(w, '-', err.Error())
	} else {
		msg := fmt.Sprintf("%v was deleted successfully", tableName)
		OutputMessage(w, '+', msg)
	}
}

/*
	List all tables, created by the user.
*/

func listTables(r *db.SQLiteRepository, w io.Writer) {
	t := r.ListTables()
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Tables:")
	for _, tt := range t {
		fmt.Fprintf(w, "     %v\n", tt)
	}

	fmt.Fprintln(w)
}

/*
	Import csv entries into the current table.
*/

func importCSV(r *db.SQLiteRepository, rd io.Reader, w io.Writer) {
	entries, err := importer.ImportCSV(rd)
	if err != nil {
		OutputMessage(w, '-', err.Error())
		return
	}

	for i, e := range entries {
		_, err = r.Insert(*e)
		if err != nil {
			if errors.Is(err, db.ErrDuplicate) {
				msg := fmt.Sprintf("Entry on line %d already exists", i+2)
				OutputMessage(w, '-', msg)
				return
			}

			OutputMessage(w, '-', err.Error())
			return
		}
	}

	msg := fmt.Sprintf("import completed successfully. %d entries added.",
		len(entries))
	OutputMessage(w, '+', msg)
}

/*
	Converts the entries within the current table to XML and write it to the
	out io.Writer
*/

func exportTable(r *db.SQLiteRepository, w, out io.Writer) {
	entries, err := r.All()

	if err != nil {
		OutputMessage(w, '-', err.Error())
		return
	}

	book, err := exporter.ExportAddressBook(entries)
	if err != nil {
		OutputMessage(w, '-', err.Error())
		return
	}

	s := exporter.ElementToString(book)
	out.Write([]byte(s))

	msg := "table exported successfully"
	OutputMessage(w, '+', msg)
}

/*
	Writes a help message to w, typically os.Stdout
*/

func helpUser(w io.Writer) {
	msg := "type 'help' for a list of commands"
	OutputMessage(w, '!', msg)
}
