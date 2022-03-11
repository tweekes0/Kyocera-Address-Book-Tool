package prompt

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/kitar0s/kyocera-ab-tool/db"
)

func TestPrompt(t *testing.T) {
	repo, teardown := db.SetupWithInserts(t)
	defer teardown()

	tt := []struct {
		description string
		input string
		expected string
	} {
		{
			description: "add_user command",
			input: "add_user john doe,jdoe,jdoe@email.com",
			expected: "[+] john doe was added successfully\n\n",
		},
		{
			description: "add_user command for duplicate",
			input: "add_user john doe,jdoe,jdoe@email.com",
			expected: "[-] record already exists\n\n",
		},
		{
			description: "delete_user command",
			input: "delete_user jdoe",
			expected: "[+] john doe was deleted successfully\n\n",
		},
		{
			description: "delete_user command without param",
			input: "delete_user",
			expected: "\ndelete a single user from the current table\nusage: delete_user 'USERNAME'\n\n",
		},
		{
			description: "unknown command",
			input: "unknown",
			expected: "[!] type 'help' for a list of commands\n\n",
		},
		{
			description: "blank link",
			input: "",
			expected: "",
		},
	}

	for _, tc := range tt  {
		t.Run(tc.description, func(t *testing.T) {	
			var got bytes.Buffer
			rd := ioutil.NopCloser(strings.NewReader(tc.input))
			Prompt(repo, rd, &got)

			if got.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", got.String(), tc.expected)
			}
		})
	}
}

func TestOutputMessage(t *testing.T) {
	tt := [] struct {
		description string
		input struct {
			symbol rune
			msg string
		}
		expected string
	} {
		{
			description: "output error message",
			input: struct{symbol rune; msg string}{
				symbol: '-',
				msg: "ERROR!",
			},
			expected: "[-] ERROR!\n\n",
		},
		{
			description: "output success message",
			input: struct{symbol rune; msg string}{
				symbol: '+',
				msg: "SUCCESS!",
			},
			expected: "[+] SUCCESS!\n\n",
		},
		{
			description: "output exclamatory message",
			input: struct{symbol rune; msg string}{
				symbol: '!',
				msg: "EXPLANATION!",
			},
			expected: "[!] EXPLANATION!\n\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			var got bytes.Buffer
			outputMessage(&got, tc.input.symbol, tc.input.msg)

			if got.String() != tc.expected {
				t.Fatalf("got: %v, expected: %v", got.String(), tc.expected)
			}
		})
	}
}

func TestStripQuotes(t *testing.T) {
	tt := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "remove single quotes from string",
			input:       `I 'want' single-quotes 'GONE'`,
			expected:    "I want single-quotes GONE",
		},
		{
			description: "remove double quotes from string",
			input:       `I "hate" double-quotes`,
			expected:    "I hate double-quotes",
		},
		{
			description: "remove all types quotes and back ticks",
			input:       "`````h`\"\"\"'''''''i",
			expected:    "hi",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			got := stripQuotes(tc.input)

			if got != tc.expected {
				t.Fatalf("got: %v, expected: %v", got, tc.expected)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	tt := []struct {
		description string
		input       string
		expected    struct {
			command string
			param   string
		}
	}{
		{
			description: "parse create_table commands",
			input:       "create_table new_table",
			expected: struct {
				command string
				param   string
			}{
				command: "create_table",
				param:   "new_table",
			},
		},
		{
			description: "parse import_csv command with file path",
			input:       "import_csv /path/to/a/really really/cool/file",
			expected: struct {
				command string
				param   string
			}{
				command: "import_csv",
				param:   "/path/to/a/really really/cool/file",
			},
		},
		{
			description: "parse single command with many spaces",
			input:       "                help         ",
			expected: struct {
				command string
				param   string
			}{
				command: "help",
				param:   "",
			},
		},
		{
			description: "parse command and param with many spaces",
			input:       "          create_table    another_table       ",
			expected: struct {
				command string
				param   string
			}{
				command: "create_table",
				param:   "another_table",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			command, param := parseArgs(tc.input)

			if command != tc.expected.command {
				t.Fatalf("got: %v, expected: %v", command, tc.expected.command)
			}

			if param != tc.expected.param {
				t.Fatalf("got: %v, expected: %v", param, tc.expected.param)
			}
		})
	}
}
