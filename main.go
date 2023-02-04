package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

func main() {

	wss := NewWSServer("0.0.0.0", 9696, func(clientID string, messageType int, message []byte) {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("[%s] %s", red(clientID), string(message))
	})

	wss.Start()

	shell := NewShell(func(tokens []string) {
		//println(strings.Join(tokens, " "))
		clientID := tokens[0]
		wss.WriteText(clientID, strings.Join(tokens[1:], " "))
	})

	err := shell.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

}
