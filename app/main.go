package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/commands"
)

func parseLine(line string) (command string, args []string) {
	words := strings.Split(strings.TrimSpace(line), " ")
	command = words[0]
	for _, arg := range words[1:] {
		if arg == "" {
			continue
		}
		args = append(args, arg)
	}
	return
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("$ ")
		if scanner.Scan() {
			command, args := parseLine(scanner.Text())

			err := commands.ExecuteBuiltin(command, args)
			if errors.Is(err, commands.ErrNotBuiltin) {
				// Not a builtin command, check the path instead
				cmd := exec.Command(command, args...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				err := cmd.Run()
				if err != nil {
					if errors.Is(err, exec.ErrNotFound) {
						fmt.Printf("%s: command not found\n", command)
					} else {
						fmt.Println(err)
					}
				}
			}
		}
	}
}
