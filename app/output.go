package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func getOutputFile(path string, append bool) (*os.File, bool) {
	flag := os.O_CREATE | os.O_WRONLY
	if append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	if !filepath.IsAbs(path) {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return nil, false
		}
		path = filepath.Join(pwd, path)
	}

	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return nil, false
	}
	return file, true
}

func getStdoutFile(path string, append bool) (*os.File, bool) {
	if path == "" {
		return os.Stdout, false
	}
	return getOutputFile(path, append)
}

func getStderrFile(path string, append bool) (*os.File, bool) {
	if path == "" {
		return os.Stderr, false
	}
	return getOutputFile(path, append)
}
