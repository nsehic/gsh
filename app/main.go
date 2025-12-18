package main

import (
	"bufio"
	"fmt"
	"os"
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
		fmt.Fprint(os.Stdout, "$ ")
		if scanner.Scan() {
			command, args := parseLine(scanner.Text())

			if command == "" {
				continue
			} else if commands.IsBuiltin(command) {
				commands.Builtins[command](args)
			} else {
				fmt.Printf("%s: command not found\n", command)
			}
		}
	}
}
