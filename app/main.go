package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codecrafters-io/shell-starter-go/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/commands"
)

func findExecutablesInPath() ([]string, error) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil, nil
	}

	seen := make(map[string]struct{})
	var executables []string

	for _, dir := range filepath.SplitList(pathEnv) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			if info.Mode().IsRegular() && info.Mode()&0111 != 0 {
				name := entry.Name()
				if _, ok := seen[name]; !ok {
					seen[name] = struct{}{}
					executables = append(executables, name)
				}
			}
		}
	}

	return executables, nil
}

func main() {
	allCommands := []string{"echo", "exit"}
	pathExecutables, _ := findExecutablesInPath()
	allCommands = append(allCommands, pathExecutables...)

	lr := LineReader{
		autocompletion: autocomplete.NewTrie(allCommands...),
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
