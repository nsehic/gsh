package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtins = map[string]struct{}{
	"echo": {},
	"exit": {},
	"type": {},
}

func isExecutable(entry os.DirEntry) bool {
	info, err := entry.Info()
	if err != nil {
		return false
	}

	mode := info.Mode()
	return mode.IsRegular() && mode&0111 != 0
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

	if _, ok := builtins[command]; ok {
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
