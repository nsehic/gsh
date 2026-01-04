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

var ErrBuiltinNotExists = errors.New("this builtin does not exist")

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
		return ErrBuiltinNotExists
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
		fmt.Println("type: invalid arguments")
		return
	}

	for _, arg := range args {
		if IsBuiltin(arg) {
			fmt.Printf("%s is a shell builtin\n", arg)
			continue
		}

		path, err := exec.LookPath(arg)
		if err == nil {
			fmt.Printf("%s is %s\n", arg, path)
			continue
		}

		fmt.Printf("%s: not found\n", arg)
	}
}

func pwdCmd(args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("pwd: %v\n", err)
	}
	fmt.Println(dir)
}

func cdCmd(args []string) {
	if len(args) > 1 {
		fmt.Println("Too many args to cd command")
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("cd: %v\n", err)
	}
	var dirpath string
	if len(args) == 0 {
		dirpath = "~"
	} else {
		dirpath = args[0]
	}
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
