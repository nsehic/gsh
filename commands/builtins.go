package commands

import (
	"fmt"
	"strings"
)

var builtins = map[string]bool{
	"echo": true,
	"exit": true,
	"type": true,
}

func Echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func Type(args []string) {
	if len(args) < 1 {
		fmt.Println("Provide a valid command")
		return
	}

	command := args[0]
	if builtin := builtins[command]; builtin {
		fmt.Printf("%s is a shell builtin\n", command)
	} else {
		fmt.Printf("%s: not found\n", command)
	}
}
