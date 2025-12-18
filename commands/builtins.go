package commands

import (
	"fmt"
	"os"
	"strings"
)

type Command func(args []string)

var Builtins = map[string]Command{}

func register(command string, fn Command) {
	Builtins[command] = fn
}

func init() {
	register("exit", Exit)
	register("echo", Echo)
	register("type", Type)
}

func IsBuiltin(name string) bool {
	_, ok := Builtins[name]
	return ok
}

func Echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func Exit(args []string) {
	os.Exit(0)
}

func Type(args []string) {
	if len(args) < 1 {
		fmt.Println("Provide a valid command")
		return
	}

	command := args[0]

	if IsBuiltin(command) {
		fmt.Printf("%s is a shell builtin\n", command)
	} else {
		fmt.Printf("%s: not found\n", command)
	}
}
