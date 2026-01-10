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

func getFile(path string) (*os.File, error) {
	if !filepath.IsAbs(path) {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(pwd, path)
	}
	return os.Create(path)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var parser Parser

	for {
		fmt.Print("$ ")
		if scanner.Scan() {
			parseResult := parser.Parse(scanner.Text())
			if parseResult == nil {
				continue
			}

			stdout := os.Stdout
			shouldCloseStdout := false

			stderr := os.Stderr
			shouldCloseStderr := false

			if parseResult.StdoutRedirectPath != "" {
				file, err := getFile(parseResult.StdoutRedirectPath)
				if err != nil {
					fmt.Printf("error: %v\n", err)
				} else {
					stdout = file
					shouldCloseStdout = true
				}
			}

			if parseResult.StderrRedirectPath != "" {
				file, err := getFile(parseResult.StderrRedirectPath)
				if err != nil {
					fmt.Printf("error: %v\n", err)
				} else {
					stderr = file
					shouldCloseStderr = true
				}
			}

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
