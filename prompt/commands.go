package prompt

import (
	"fmt"
	"io"

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
	readline.PcItem("show_users"),
	readline.PcItem("add_user"),
	readline.PcItem("delete_user"),
)

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
			fmt.Fprintf(w, "[-] %q\n\n", err)
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
		fmt.Fprintf(w, "[-] %q\n\n", err)
	}
	err = r.Delete(username)
	if err != nil {
		fmt.Fprintf(w, "[-] %q\n\n", err)
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
		fmt.Fprintf(w, "[-] %q\n\n", err)
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
		fmt.Fprintf(w, "[-] %q\n\n", err)
	} else {
		fmt.Fprintf(w, "[+] %v was deleted sucessfully\n\n",
			tableName)
	}
}