package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stdout, "$ ")
repl:
	for scanner.Scan() {
		command := strings.TrimSpace(scanner.Text())
		switch command {
		case "exit":
			break repl
		}
		fmt.Printf("%s: command not found\n", command)
		fmt.Fprint(os.Stdout, "$ ")
	}
}
