package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/commands"
)

func getNextChar(idx int, str string) string {
	r := []rune(str)
	if idx >= len(str)-1 {
		return ""
	}
	return string(r[idx+1])
}

func parseInput(input string) (string, []string) {
	line := []string{}
	singleQuote := false
	concatString := false
	var sb strings.Builder
	for i, c := range input {
		switch string(c) {
		case "'":
			if singleQuote {
				if concatString {
					concatString = false
					continue
				} else if getNextChar(i, input) == "'" {
					concatString = true
				} else {
					singleQuote = false
					if sb.Len() > 0 {
						line = append(line, sb.String())
						sb.Reset()
					}
				}
			} else {
				singleQuote = true
			}
		case " ":
			if singleQuote {
				sb.WriteRune(c)
			} else if sb.Len() > 0 {
				line = append(line, sb.String())
				sb.Reset()
			}
		default:
			sb.WriteRune(c)
		}
	}

	if sb.Len() > 0 {
		line = append(line, sb.String())
		sb.Reset()
	}

	return line[0], line[1:]
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("$ ")
		if scanner.Scan() {
			command, args := parseInput(scanner.Text())

			err := commands.ExecuteBuiltin(command, args)
			if errors.Is(err, commands.ErrBuiltinNotExists) {
				// Not a builtin command, check the path instead
				cmd := exec.Command(command, args...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				err := cmd.Run()
				if err != nil {
					if errors.Is(err, exec.ErrNotFound) {
						fmt.Printf("%s: command not found\n", command)
					} else {
						fmt.Println(err)
					}
				}
			}
		}
	}
}
