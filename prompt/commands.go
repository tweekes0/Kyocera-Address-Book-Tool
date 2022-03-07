package prompt

import (
	"fmt"
	"io"

	"github.com/chzyer/readline"
	"github.com/kitar0s/kyocera-ab-tool/db"
)

var completions = readline.NewPrefixCompleter(
	readline.PcItem("create_table"),
	readline.PcItem("switch_table"),
	readline.PcItem("clear_table"),
	readline.PcItem("delete_table"),
	readline.PcItem("list_tables"),
	readline.PcItem("users"),
	readline.PcItem("delete_user"),
	
)

func createTable(r *db.SQLiteRepository, w io.Writer, tableName string) {
	err := r.NewTable(tableName)
	if err != nil {
		fmt.Fprintf(w, "[-] %v\n\n", err)
	} else {
		fmt.Fprintf(w, "[+] %v was created successfully\n\n",
			r.CurrentTable())
	}
}

func switchTable(r *db.SQLiteRepository, w io.Writer, tableName string) {
	err := r.SwitchTable(tableName)
	if err != nil {
		fmt.Fprintf(w, "[-] %v\n\n", err)
	}
}

func showUsers(r *db.SQLiteRepository, w io.Writer) {
	all, err := r.All()
	if err != nil {
		fmt.Fprintf(w, "[-] %v\n", err)
	} else {
		// TODO: implement pretty way to print all the entries
		fmt.Fprintf(w, "[+] contents of %v\n", r.CurrentTable())
		for _, e := range all {
			// fmt.Fprintf(w, "%v\n", *e)
			e.Display(w)
		}
	}
}

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