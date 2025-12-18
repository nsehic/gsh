package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

func main() {
	fmt.Fprint(os.Stdout, "$ ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s: command not found\n", line)
	}
}
