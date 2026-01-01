package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Command func(args []string)

var ErrNotBuiltin = errors.New("the given command is not a builtin")

var builtins map[string]Command

func init() {
	builtins = map[string]Command{
		"echo": echoCmd,
		"exit": exitCmd,
		"type": typeCmd,
	}
}

func IsBuiltin(command string) bool {
	_, ok := builtins[command]
	return ok
}

func ExecuteBuiltin(command string, args []string) error {
	if !IsBuiltin(command) {
		return ErrNotBuiltin
	}

	builtins[command](args)
	return nil
}

func echoCmd(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func exitCmd(args []string) {
	os.Exit(0)
}

func typeCmd(args []string) {
	if len(args) < 1 {
		fmt.Println("Provide a valid command")
		return
	}

	command := args[0]

	if IsBuiltin(command) {
		fmt.Printf("%s is a shell builtin\n", command)
		return
	}

	path, err := exec.LookPath(command)
	if err == nil {
		fmt.Printf("%s is %s\n", command, path)
		return
	}

	fmt.Printf("%s: not found\n", command)
}
