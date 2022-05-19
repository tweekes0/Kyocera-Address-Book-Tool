package prompt

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

/*
	Returns customized readline instance.
*/

func newReadLine(rd io.ReadCloser) *readline.Instance {
	l, err := readline.NewEx(&readline.Config{
		Prompt:            "",
		AutoComplete:      completions,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		Stdin:             rd,
	})
	if err != nil {
		log.Fatalf("could not create readline: %q", err)
	}

	return l
}

/*
	Outputs message to the user.
		+ indicates success
		- indicates error
		! indicates exclamatory
*/

func OutputMessage(w io.Writer, symbol rune, msg string) {
	fmt.Fprintf(w, "[%v] %v\n\n", string(symbol), msg)
}

/*
	Remove the quotes from the supplied string
*/

func stripQuotes(s string) string {
	s = strings.Replace(s, `"`, "", -1)
	s = strings.Replace(s, `'`, "", -1)
	s = strings.Replace(s, "`", "", -1)

	return s
}

/*
	Take a from user input and returns two strings a command and it's optional
	parameter
*/

func parseArgs(s string) (command string, param string) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		command = ""
		param = ""
		return
	}

	if len(fields) == 1 {
		command = fields[0]
		param = ""
		return
	}

	command = fields[0]
	param = strings.Join(fields[1:], " ")
	param = stripQuotes(param)

	return
}

func createFile(tblName string) (*os.File, error) {
	_, err := os.Stat("./Address Book")
	fname := fmt.Sprintf("./Address Books/%v %s.xml",
		tblName, time.Now().Format("2006-Jan-02"))

	if os.IsNotExist(err) {
		if err = os.Mkdir("./Address Books", os.ModePerm); err != nil {
			return nil, err
		}
	}

	f, err := os.Create(fname)
	if err != nil {
		return nil, err
	}

	return f, nil
}
