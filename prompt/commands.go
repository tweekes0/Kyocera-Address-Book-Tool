package prompt

import (
	"fmt"
	"io"
	"sort"

	"github.com/chzyer/readline"
	"github.com/kitar0s/kyocera-ab-tool/db"
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

	readline.PcItem("help",
		readline.PcItem("create_table"),
		readline.PcItem("switch_table"),
		readline.PcItem("clear_table"),
		readline.PcItem("delete_table"),
		readline.PcItem("list_tables"),
		readline.PcItem("show_users"),
		readline.PcItem("add_user"),
		readline.PcItem("delete_user"),
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
	Lists all the commands to the user.
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
		fmt.Fprintf(w, "[-] %v\n\n", err)
	} else {
		fmt.Fprintf(w, "[+] %v was created successfully\n\n",
			r.CurrentTable())
	}
}

/*
	Switch to a table and write to w if the operation fails.
*/

func switchTable(r *db.SQLiteRepository, w io.Writer, tableName string) {
	err := r.SwitchTable(tableName)
	if err != nil {
		fmt.Fprintf(w, "[-] %v\n\n", err)
	}
}

/*
	Display all the users that in the current table.
*/

func showUsers(r *db.SQLiteRepository, w io.Writer) {
	all, err := r.All()
	if err != nil {
		fmt.Fprintf(w, "[-] %v\n", err)
	} else {
		// TODO: implement pretty way to print all the entries
		fmt.Fprintf(w, "[+] contents of %v\n", r.CurrentTable())
		for _, e := range all {
			e.Display(w)
		}
	}
}

/*
	Inserts a new Entry into the base, granted that the params are valid
*/

func insertEntry(r *db.SQLiteRepository, w io.Writer, params []string) {
	e, err := db.NewEntry(1, params[0], params[1], params[2])
	if err != nil {
		fmt.Fprintln(w, "[-] ", err)
	} else {
		_, err = r.Insert(*e)
		if err != nil {
			fmt.Fprintf(w, "[-] %v\n\n", err)
		} else {
			fmt.Fprintf(w, "[+] %v was added successfully\n\n",
				e.Name)
		}
	}
}

/*
	Deletes an Entry from the database given a valid username.
*/

func deleteEntry(r *db.SQLiteRepository, w io.Writer, username string) {
	e, err := r.GetByUsername(username)
	if err != nil {
		fmt.Fprintf(w, "[-] %v\n\n", err)
	}
	err = r.Delete(username)
	if err != nil {
		fmt.Fprintf(w, "[-] %v\n\n", err)
	} else {
		fmt.Fprintf(w, "[+] %v was deleted sucessfully\n\n", e.Name)
	}
}

/*
	Clears all the entries from the current table
*/

func clearTable(r *db.SQLiteRepository, w io.Writer) {
	err := r.ClearTable()
	if err != nil {
		fmt.Fprintf(w, "[-] %v\n\n", err)
	} else {
		fmt.Fprintf(w, "[+] %v was cleared sucessfully\n\n",
			r.CurrentTable())
	}
}

/*
	Deletes the specified table. DEFAULT_TABLE cannot be deleted.
*/

func deleteTable(r *db.SQLiteRepository, w io.Writer, tableName string) {
	err := r.DeleteTable(tableName)
	if err != nil {
		fmt.Fprintf(w, "[-] %v\n\n", err)
	} else {
		fmt.Fprintf(w, "[+] %v was deleted sucessfully\n\n",
			tableName)
	}
}

/*
	Lists all tables, created by the user.
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
