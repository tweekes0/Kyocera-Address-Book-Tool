package prompt

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/chzyer/readline"
	"github.com/kitar0s/kyocera-ab-tool/db"
	"github.com/kitar0s/kyocera-ab-tool/importer"
)

/*
	List of PcItems that readline structs use for terminal autocompletes.
*/

var completions = readline.NewPrefixCompleter(
	readline.PcItem("create_table"),
	readline.PcItem("switch_table"),
	readline.PcItem("clear_table"),
	readline.PcItem("delete_table"),
	readline.PcItem("list_tables"),
	readline.PcItem("show_users"),
	readline.PcItem("add_user"),
	readline.PcItem("delete_user"),
	readline.PcItem("import_csv"),

	readline.PcItem("help",
		readline.PcItem("create_table"),
		readline.PcItem("switch_table"),
		readline.PcItem("clear_table"),
		readline.PcItem("delete_table"),
		readline.PcItem("list_tables"),
		readline.PcItem("show_users"),
		readline.PcItem("add_user"),
		readline.PcItem("delete_user"),
		readline.PcItem("import_csv"),
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
		description: "clear the current table of all entries",
		usage:       "clear_table",
	},
	"list_tables": {
		description: "list all tables",
		usage:       "list_table",
	},
	"show_users": {
		description: "show all the users in the current table",
		usage:       "show_users",
	},
	"add_user": {
		description: "add user to the current table. Fields must be separated by commas",
		usage:       "add_user NAME,USERNAME,EMAIL",
	},
	"delete_user": {
		description: "delete a single user from the current table",
		usage:       "delete_user USERNAME",
	},
	"import_csv": {
		description: "import users from csv file into current table",
		usage:       "import_csv PATH_TO_FILE",
	},
}

/*
	Function that will print the description and usage of a commmand if it exists,
	otherwise it will list the available commands.
*/

func helpCommand(w io.Writer, s string) {
	if command, ok := commands[s]; ok {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "%v\n", command.description)
		fmt.Fprintf(w, "usage: %v\n", command.usage)
		fmt.Fprintln(w)
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

	fmt.Fprintln(w)
	fmt.Fprint(w, "Commands:\n")
	for _, k := range keys {
		fmt.Fprintf(w, "%-15v : %10v\n", k, commands[k].description)
	}

	fmt.Fprintln(w)
}

/*
	Create a new table and write to io.Writer the success or failure
	of the operation.
*/

func createTable(r *db.SQLiteRepository, w io.Writer, tableName string) {
	err := r.NewTable(tableName)
	if err != nil {
		outputMessage(w, '-', err.Error())
	} else {
		msg := fmt.Sprintf("%v was created successfully", r.CurrentTable())
		outputMessage(w, '+', msg)
	}
}

/*
	Switch to a table and write to w if the operation fails.
*/

func switchTable(r *db.SQLiteRepository, w io.Writer, tableName string) {
	err := r.SwitchTable(tableName)
	if err != nil {
		outputMessage(w, '-', err.Error())
	}
}

/*
	Display all the users that in the current table.
*/

func showUsers(r *db.SQLiteRepository, w io.Writer) {
	all, err := r.All()
	if err != nil {
		outputMessage(w, '-', err.Error())
	} else {
		// TODO: implement pretty way to print all the entries
		msg := fmt.Sprintf("contents of %v", r.CurrentTable())
		outputMessage(w, '+', msg)
		for _, e := range all {
			e.Display(w)
		}
	}
}

/*
	Inserts a new Entry into the base, granted that the params are valid
*/

func insertEntry(r *db.SQLiteRepository, w io.Writer, params []string) {
	e, err := db.NewEntry(params[0], params[1], params[2])
	if err != nil {
		outputMessage(w, '-', err.Error())
	} else {
		_, err = r.Insert(*e)
		if err != nil {
			outputMessage(w, '-', err.Error())
		} else {
			msg := fmt.Sprintf("%v was added successfully", e.Name)
			outputMessage(w, '+', msg)
		}
	}
}

/*
	Delete an Entry from the database given a valid username.
*/

func deleteEntry(r *db.SQLiteRepository, w io.Writer, username string) {
	e, err := r.GetByUsername(username)
	if err != nil {
		outputMessage(w, '-', err.Error())
	}
	err = r.Delete(username)
	if err != nil {
		outputMessage(w, '-', err.Error())
	} else {
		msg := fmt.Sprintf("%v was deleted sucessfully", e.Name)
		outputMessage(w, '+', msg)
	}
}

/*
	Clear all the entries from the current table
*/

func clearTable(r *db.SQLiteRepository, w io.Writer) {
	err := r.ClearTable()
	if err != nil {
		outputMessage(w, '-', err.Error())
	} else {
		msg := fmt.Sprintf("%v was cleared sucessfully", r.CurrentTable())
		outputMessage(w, '+', msg)

	}
}

/*
	Deletes the specified table. DEFAULT_TABLE cannot be deleted.
*/

func deleteTable(r *db.SQLiteRepository, w io.Writer, tableName string) {
	err := r.DeleteTable(tableName)
	if err != nil {
		outputMessage(w, '-', err.Error())
	} else {
		msg := fmt.Sprintf("%v was deleted successfully", tableName)
		outputMessage(w, '+', msg)
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
	Import csv entries into the current table
*/

func importCSV(r *db.SQLiteRepository, rd io.Reader, w io.Writer) {
	entries, err := importer.ImportCSV(rd)
	if err != nil {
		outputMessage(w, '-', err.Error())
		return
	}

	for i, e := range entries {
		_, err = r.Insert(*e)
		if err != nil {
			if errors.Is(err, db.ErrDuplicate) {
				msg := fmt.Sprintf("Entry on line %d already exists", i+2)
				outputMessage(w, '-', msg)
				return
			}

			outputMessage(w, '-', err.Error())
			return
		}
	}

	msg := fmt.Sprintf("import completed successfully. %d entries added.",
		len(entries))
	outputMessage(w, '+', msg)
}
