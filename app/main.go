package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/codecrafters-io/shell-starter-go/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/commands"
	"golang.org/x/term"
)

const prompt = "$ "

func redrawLine(line []rune) {
	fmt.Fprint(os.Stdout, "\r\033[K")
	fmt.Fprint(os.Stdout, prompt)
	fmt.Fprint(os.Stdout, string(line))
}

func printPrompt() {
	fmt.Fprintf(os.Stdout, "\r%s", prompt)
}

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

func main() {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		fmt.Fprintf(os.Stderr, "error: stdin is not a terminal\n")
		return
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		return
	}
	defer term.Restore(fd, oldState)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		term.Restore(fd, oldState)
		os.Exit(1)
	}()

	trie := autocomplete.NewTrie("echo", "exit")

	buf := make([]byte, 1)
	line := []rune{}
	printPrompt()
	var parser Parser

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		b := buf[0]

		switch b {
		case '\r', '\n':
			os.Stdout.Write([]byte("\r\n"))

			input := string(line)
			line = nil // reset buffer

			if strings.TrimSpace(input) == "" {
				printPrompt()
				break
			}

			parseResult := parser.Parse(input)
			if parseResult != nil {
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
			printPrompt()

		case '\t':
			// Autocomplete
			completions := trie.Complete(string(line))
			if len(completions) == 0 {
				break
			}
			line = []rune(fmt.Sprintf("%s ", completions[0]))
			redrawLine(line)
		case 127:
			// Backspace
			if len(line) > 0 {
				line = line[:len(line)-1]
				os.Stdout.Write([]byte("\b \b"))
			}
		default:
			// Normal printable character
			line = append(line, rune(b))
			os.Stdout.Write([]byte{b})
		}
	}
}
