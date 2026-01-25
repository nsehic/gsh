package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/commands"
	"golang.org/x/term"
)

func getOutputFile(path string, append bool) (*os.File, bool) {
	flag := os.O_CREATE | os.O_WRONLY
	if append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	if !filepath.IsAbs(path) {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return nil, false
		}
		path = filepath.Join(pwd, path)
	}

	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return nil, false
	}
	return file, true
}

func getStdoutFile(path string, append bool) (*os.File, bool) {
	if path == "" {
		return os.Stdout, false
	}
	return getOutputFile(path, append)
}

func getStderrFile(path string, append bool) (*os.File, bool) {
	if path == "" {
		return os.Stderr, false
	}
	return getOutputFile(path, append)
}

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

func main() {
	lr := LineReader{
		autocompletion: autocomplete.NewTrie("echo", "exit"),
		prompt:         "$ ",
	}

	var parser Parser

	for {
		input, err := lr.ReadLine()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		if parseResult := parser.Parse(input); parseResult != nil {
			stdout, shouldCloseStdout := getStdoutFile(
				parseResult.StdoutPath,
				parseResult.StdoutAppend,
			)
			stderr, shouldCloseStderr := getStderrFile(
				parseResult.StderrPath,
				parseResult.StderrAppend,
			)

			cmd := commands.Command(parseResult.Command, parseResult.Args...)
			cmd.Out = stdout
			cmd.Err = stderr

			if err := cmd.Run(); err != nil {
				var execError *exec.Error
				var pathError *os.PathError

				switch {
				case errors.As(err, &execError):
					fmt.Fprintf(cmd.Err, "%s: not found\n", cmd.Name)
				case errors.As(err, &pathError):
					fmt.Fprintf(cmd.Err, "%s: %v\n", cmd.Name, pathError.Err)
				}
			}

			if shouldCloseStdout {
				stdout.Close()
			}
			if shouldCloseStderr {
				stderr.Close()
			}
		}
	}
}
