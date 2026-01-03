package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/commands"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var parser Parser

	for {
		fmt.Print("$ ")
		if scanner.Scan() {
			command, args := parser.Parse(scanner.Text())

			err := commands.ExecuteBuiltin(command, args)
			if errors.Is(err, commands.ErrBuiltinNotExists) {
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
