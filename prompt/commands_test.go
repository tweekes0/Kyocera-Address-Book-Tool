package prompt

import (
	"bytes"
	"fmt"
	"sort"
	"testing"

	"github.com/tweekes0/kyocera-ab-tool/db"
	"github.com/tweekes0/kyocera-ab-tool/importer"
)

func TestHelpCommand(t *testing.T) {
	tt := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "help command for create_table",
			input:       "create_table",
			expected:    "\ncreates new table and sets it to the current table\nusage: create_table 'TABLE_NAME'\n\n",
		},
		{
			description: "help command for switch_table",
			input:       "switch_table",
			expected:    "\nswitch the current table\nusage: switch_table 'TABLE_NAME'\n\n",
		},
		{
			description: "help command for list_tables",
			input:       "list_tables",
			expected:    "\nlist all tables\nusage: list_tables\n\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			var got bytes.Buffer
			helpCommand(&got, tc.input)

			if got.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", got.String(), tc.expected)
			}
		})
	}

	t.Run("help command without a paramter", func(t *testing.T) {
		t.Parallel()

		keys := make([]string, 0, len(commands))
		for k := range commands {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var buf bytes.Buffer
		var expected bytes.Buffer

		helpCommand(&buf, "")

		expected.WriteString("\nCommands:\n")

		for _, k := range keys {
			expected.WriteString(fmt.Sprintf("     %-15v : %10v\n",
				k, commands[k].description))
		}
		expected.WriteString("\n")

		if buf.String() != expected.String() {
			t.Fatalf("got: %v, expected: %v", buf.String(), expected.String())
		}
	})
}

func TestListCommands(t *testing.T) {
	t.Parallel()

	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var got bytes.Buffer
	var expected bytes.Buffer

	listCommands(&got)

	expected.WriteString("\nCommands:\n")

	for _, k := range keys {
		expected.WriteString(fmt.Sprintf("     %-15v : %10v\n",
			k, commands[k].description))
	}
	expected.WriteString("\n")

	if got.String() != expected.String() {
		t.Fatalf("got: %v, expected: %v", got.String(), expected.String())
	}
}

func TestCreateTable(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	tt := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "create a valid table",
			input:       "valid_table",
			expected:    "[+] valid_table was created successfully\n\n",
		},
		{
			description: "create an existing table",
			input:       db.DEFAULT_TABLE,
			expected:    "[-] table already exists\n\n",
		},
		{
			description: "create an invalid table",
			input:       "--invalid--",
			expected:    "[-] tablename is not valid\n\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {

			var got bytes.Buffer
			createTable(repo, &got, tc.input)

			if got.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", got.String(), tc.expected)
			}
		})
	}
}

func TestSwitchTable(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	tt := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "switch to an existing table",
			input:       db.DEFAULT_TABLE,
			expected:    "",
		},
		{
			description: "switch to an non-existing table",
			input:       "non_existing_table",
			expected:    "[-] table does not exist\n\n",
		},
		{
			description: "switch to an invalid table",
			input:       "--invalid--",
			expected:    "[-] tablename is not valid\n\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {

			var got bytes.Buffer
			switchTable(repo, &got, tc.input)

			if got.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", got.String(), tc.expected)
			}
		})
	}
}

func TestShowUser(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	t.Run("show users in the default table after inserts", func(t *testing.T) {
		var got, expected bytes.Buffer
		all, err := repo.All()
		if err != nil {
			t.Fatal(err)
		}

		showUsers(repo, &got)

		expected.WriteString(fmt.Sprintf("[+] contents of %v\n\n",
			repo.CurrentTable()))

		for _, e := range all {
			e.Display(&expected)
		}

		if got.String() != expected.String() {
			t.Fatalf("got: %v, expected: %v", got.String(), expected.String())
		}
	})

	t.Run("show users in an empty table", func(t *testing.T) {
		var got, expected bytes.Buffer

		repo.NewTable("valid_table")
		showUsers(repo, &got)

		expected.WriteString(fmt.Sprintf("[!] %v is empty\n\n",
			repo.CurrentTable()))

		if got.String() != expected.String() {
			t.Fatalf("got: %v, expected: %v", got.String(), expected.String())
		}
	})
}

func TestAddUser(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	tt := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "add valid user",
			input:       "jane doe,jdoe,jdoe@email.com",
			expected:    "[+] jane doe was added successfully\n\n",
		},
		{
			description: "add existing user",
			input:       "jane doe,jdoe,jdoe@email.com",
			expected:    "[-] record already exists\n\n",
		},
		{
			description: "add invalid user",
			input:       "jane 1,jdoe,jdoe@email.com",
			expected:    "[-] name is not valid\n\n",
		},
		{
			description: "add user with too few fields ",
			input:       "jane 1,jdoe",
			expected:    "[-] invalid number of fields\n\n",
		},
	}

	for _, tc := range tt {
		var got bytes.Buffer

		t.Run(tc.description, func(t *testing.T) {
			addUser(repo, &got, tc.input)

			if got.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", got.String(), tc.expected)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	tt := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "delete a valid user",
			input:       "username1",
			expected:    "[+] Test One was deleted successfully\n\n",
		},
		{
			description: "delete a non-existing user",
			input:       "username1",
			expected:    "[-] record does not exist\n\n",
		},
		{
			description: "delete an invalid user",
			input:       "--invalid--",
			expected:    "[-] username is not valid\n\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			var got bytes.Buffer
			deleteUser(repo, &got, tc.input)

			if got.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", got.String(), tc.expected)
			}
		})
	}
}

func TestClearTable(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	var got, expected bytes.Buffer
	expected.WriteString(fmt.Sprintf("[+] %v was cleared successfully\n\n",
		repo.CurrentTable()))
	clearTable(repo, &got)

	all, _ := repo.All()

	if got.String() != expected.String() {
		t.Fatalf("got: %v, expected: %v", got.String(), expected.String())
	}

	if len(all) != 0 {
		t.Fatalf("got: %v, expected: %v", all, []*db.Entry{})
	}
}

func TestDeleteTable(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	repo.NewTable("valid")

	tt := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "delete default table",
			input:       db.DEFAULT_TABLE,
			expected:    "[-] table cannot be deleted\n\n",
		},
		{
			description: "delete valid table",
			input:       "valid",
			expected:    "[+] valid was deleted successfully\n\n",
		},
		{
			description: "delete invalid table",
			input:       "--invalid",
			expected:    "[-] tablename is not valid\n\n",
		},
	}

	for _, tc := range tt {
		var got bytes.Buffer
		deleteTable(repo, &got, tc.input)

		if got.String() != tc.expected {
			t.Fatalf("got: %v, expected: %v", got.String(), tc.expected)
		}
	}
}

func TestListTables(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	repo.NewTable("another_table")
	repo.NewTable("last_table")

	tables := []string{"default_table", "another_table", "last_table"}
	var got, expected bytes.Buffer

	expected.WriteString("\nTables:\n")
	for _, table := range tables {
		expected.WriteString(fmt.Sprintf("     %v\n", table))
	}
	expected.WriteString("\n")

	listTables(repo, &got)

	if got.String() != expected.String() {
		t.Fatalf("got: %v, expected: %v", got.String(), expected.String())
	}
}

func TestImportCSV(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	csv1 := [][]string{
		{"name", "username", "email"},
		{"Jane Doe", "janedoe", "janedoe@email.com"},
		{"John Doe", "johndoe", "johndoe@email.com"},
	}

	csv2 := [][]string{
		{"invalid", "username", "email"},
		{"valid name", "janedoe", "janedoe@email.com"},
	}

	tt := []struct {
		description string
		input [][]string
		expected string
	} {
		{
			description: "import valid csv",
			input: csv1,
			expected: "[+] import completed successfully. 2 entries added.\n\n",
		},
		{
			description: "import csv with invalid header",
			input: csv2,
			expected: "[-] invalid header\n\n",
		},
		{
			description: "import csv with existing entry",
			input: csv1,
			expected: "[-] Entry on line 2 already exists\n\n",
		},
	}

	for _, tc := range tt {

		t.Run(tc.description, func(t *testing.T) {
			r, td := importer.SetupCSV(t, tc.input)
			defer td()

			var got bytes.Buffer
			importCSV(repo, r, &got)

			if got.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", got.String(), tc.expected)
			}

		})
	}
}

func TestHelpUser(t *testing.T) {
	t.Parallel()
	
	var got bytes.Buffer
	expected := "[!] type 'help' for a list of commands\n\n"

	helpUser(&got)

	if got.String() != expected {
		t.Fatalf("got: %v, expected: %v", got.String(), expected)
	}
}
