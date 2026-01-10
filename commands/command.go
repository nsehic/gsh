package commands

import (
	"io"
	"os"
	"os/exec"
)

type Cmd struct {
	Name string
	Args []string
	In   io.Reader
	Out  io.Writer
	Err  io.Writer

	isBuiltin bool
}

func Command(name string, arg ...string) *Cmd {
	return &Cmd{
		Name:      name,
		Args:      append([]string{}, arg...),
		In:        os.Stdin,
		Out:       os.Stdout,
		Err:       os.Stderr,
		isBuiltin: isBuiltin(name),
	}
}

func (c *Cmd) Run() error {
	if c.isBuiltin {
		builtins[c.Name](c)
		return nil
	}
	cmd := exec.Command(c.Name, c.Args...)
	cmd.Stdin = c.In
	cmd.Stdout = c.Out
	cmd.Stderr = c.Err

	return cmd.Run()
}
