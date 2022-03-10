package prompt

import (
	"bytes"
	"fmt"
	"sort"
	"testing"

	"github.com/kitar0s/kyocera-ab-tool/db"
)

func TestHelpCommand(t *testing.T) {
	tt := []struct {
		description string
		got         string
		expected    string
	}{
		{
			description: "help command for create_table",
			got:         "create_table",
			expected:    "\ncreates new table and sets it to the current table\nusage: create_table 'TABLE_NAME'\n",
		},
		{
			description: "help command for switch_table",
			got:         "switch_table",
			expected:    "\nswitch the current table\nusage: switch_table 'TABLE_NAME'\n",
		},
		{
			description: "help command for list_tables",
			got:         "list_tables",
			expected:    "\nlist all tables\nusage: list_tables\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			helpCommand(&buf, tc.got)

			if buf.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", buf.String(), tc.expected)
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

	var buf bytes.Buffer
	var expected bytes.Buffer

	listCommands(&buf)

	expected.WriteString("\nCommands:\n")

	for _, k := range keys {
		expected.WriteString(fmt.Sprintf("     %-15v : %10v\n",
			k, commands[k].description))
	}
	expected.WriteString("\n")

	if buf.String() != expected.String() {
		t.Fatalf("got: %v, expected: %v", buf.String(), expected.String())
	}
}

func TestCreateTable(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	tt := []struct {
		description string
		got         string
		expected    string
	}{
		{
			description: "create a valid table",
			got:         "valid_table",
			expected:    "[+] valid_table was created successfully\n\n",
		},
		{
			description: "create an existing table",
			got:         db.DEFAULT_TABLE,
			expected:    "[-] table already exists\n\n",
		},
		{
			description: "create an invalid table",
			got:         "--invalid--",
			expected:    "[-] tablename is not valid\n\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {

			var buf bytes.Buffer
			createTable(repo, &buf, tc.got)

			if buf.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", buf.String(), tc.expected)
			}
		})
	}
}

func TestSwitchTable(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	tt := []struct {
		description string
		got         string
		expected    string
	}{
		{
			description: "switch to an existing table",
			got:         db.DEFAULT_TABLE,
			expected:    "",
		},
		{
			description: "switch to an non-existing table",
			got:         "non_existing_table",
			expected:    "[-] table does not exist\n\n",
		},
		{
			description: "switch to an invalid table",
			got:         "--invalid--",
			expected:    "[-] tablename is not valid\n\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {

			var buf bytes.Buffer
			switchTable(repo, &buf, tc.got)

			if buf.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", buf.String(), tc.expected)
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
