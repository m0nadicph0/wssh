package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	shell := NewShell(func(tokens []string) {
		println(strings.Join(tokens, " "))
	})

	err := shell.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

}
