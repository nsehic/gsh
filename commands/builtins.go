package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BuiltinCommand func(c *Cmd)

var ErrBuiltinNotExists = errors.New("this builtin does not exist")

var builtins map[string]BuiltinCommand

func init() {
	builtins = map[string]BuiltinCommand{
		"echo": echoCmd,
		"exit": exitCmd,
		"type": typeCmd,
		"pwd":  pwdCmd,
		"cd":   cdCmd,
	}
}

func isBuiltin(name string) bool {
	_, ok := builtins[name]
	return ok
}

func echoCmd(c *Cmd) {
	fmt.Fprintf(c.Out, "%s\n", strings.Join(c.Args, " "))
}

func exitCmd(c *Cmd) {
	os.Exit(0)
}

func typeCmd(c *Cmd) {
	if len(c.Args) < 1 {
		fmt.Fprint(c.Err, "type: invalid arguments\n")
		return
	}

	for _, arg := range c.Args {
		if isBuiltin(arg) {
			fmt.Fprintf(c.Out, "%s is a shell builtin\n", arg)
			continue
		}

		path, err := exec.LookPath(arg)
		if err == nil {
			fmt.Fprintf(c.Out, "%s is %s\n", arg, path)
			continue
		}

		fmt.Fprintf(c.Err, "%s: not found\n", arg)
	}
}

func pwdCmd(c *Cmd) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(c.Err, "pwd: %v\n", err)
	}
	fmt.Fprintf(c.Out, "%s\n", dir)
}

func cdCmd(c *Cmd) {
	if len(c.Args) > 1 {
		fmt.Fprint(c.Err, "Too many args to cd command\n")
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(c.Err, "cd: %v\n", err)
	}
	var dirpath string
	if len(c.Args) == 0 {
		dirpath = "~"
	} else {
		dirpath = c.Args[0]
	}
	absDirpath := dirpath

	if strings.HasPrefix(absDirpath, "~") {
		absDirpath = strings.Replace(absDirpath, "~", home, 1)
	}

	if !filepath.IsAbs(absDirpath) {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(c.Err, "cd: %v\n", err)
			return
		}
		absDirpath = filepath.Join(pwd, absDirpath)
	}

	info, err := os.Stat(absDirpath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(c.Err, "cd: %s: No such file or directory\n", dirpath)
		return
	}

	if info.IsDir() {
		if err := os.Chdir(absDirpath); err != nil {
			fmt.Fprintf(c.Err, "cd: %s: %v", dirpath, err)
		}
	}
}
