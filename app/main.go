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
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stdout, "$ ")
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s: command not found\n", line)
		fmt.Fprint(os.Stdout, "$ ")
	}
}
