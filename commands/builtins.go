package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		"pwd":  pwdCmd,
		"cd":   cdCmd,
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
	if len(args) != 1 {
		fmt.Println("type: invalid arguments")
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

func pwdCmd(args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("pwd: %v\n", err)
	}
	fmt.Println(dir)
}

func cdCmd(args []string) {
	if len(args) != 1 {
		fmt.Println("cd: invalid arguments")
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("cd: %v\n", err)
	}
	dirpath := args[0]
	absDirpath := dirpath

	if strings.HasPrefix(absDirpath, "~") {
		absDirpath = strings.Replace(absDirpath, "~", home, 1)
	}

	if !filepath.IsAbs(absDirpath) {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("cd: %v\n", err)
			return
		}
		absDirpath = filepath.Join(pwd, absDirpath)
	}

	info, err := os.Stat(absDirpath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("cd: %s: No such file or directory\n", dirpath)
		return
	}

	if info.IsDir() {
		if err := os.Chdir(absDirpath); err != nil {
			fmt.Printf("cd: %s: %v", dirpath, err)
		}
	}

}
