package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/autocomplete"
	"golang.org/x/term"
)

type LineReader struct {
	prompt         string
	autocompletion *autocomplete.Trie
}

func (lr *LineReader) RedrawLine(line string) {
	fmt.Fprint(os.Stdout, "\r\033[K")
	fmt.Fprint(os.Stdout, lr.prompt)
	fmt.Fprint(os.Stdout, line)
}

func (lr *LineReader) ReadLine() (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		fmt.Fprintf(os.Stderr, "error: stdin is not a terminal\n")
		return "", errors.New("stdin is not a terminal")
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		return "", err
	}
	defer term.Restore(fd, oldState)

	var line strings.Builder
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprintf(os.Stdout, "\r%s", lr.prompt)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return line.String(), err
		}
		switch b {
		case '\r', '\n':
			os.Stdout.Write([]byte("\r\n"))
			return line.String(), nil
		case '\t':
			completions := lr.autocompletion.Complete(line.String())
			if len(completions) == 0 {
				os.Stdout.Write([]byte("\x07"))
				break
			}
			line.Reset()
			line.WriteString(fmt.Sprintf("%s ", completions[0]))
			lr.RedrawLine(line.String())
		case 127, 8:
			if line.Len() > 0 {
				s := line.String()
				line.Reset()
				line.WriteString(s[:len(s)-1])
				os.Stdout.Write([]byte("\b \b"))
			}
		default:
			line.WriteByte(b)
			os.Stdout.Write([]byte{b})
		}
	}
}
