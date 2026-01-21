package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codecrafters-io/shell-starter-go/commands"
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	// trie := autocomplete.NewTrie("echo", "exit")
	var parser Parser

	for {
		fmt.Print("$ ")
		if scanner.Scan() {
			parseResult := parser.Parse(scanner.Text())
			if parseResult == nil {
				continue
			}

			stdout, shouldCloseStdout := getStdoutFile(parseResult.StdoutPath, parseResult.StdoutAppend)
			stderr, shouldCloseStderr := getStderrFile(parseResult.StderrPath, parseResult.StderrAppend)

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
				if err := stdout.Close(); err != nil {
					fmt.Printf("error: %v\n", err)
				}
			}

			if shouldCloseStderr {
				if err := stderr.Close(); err != nil {
					fmt.Printf("error: %v\n", err)
				}
			}
		}
	}
}
