package main

import (
	"context"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"os"
	"wssh/cmd"
	"wssh/constants"
	"wssh/wss"
)

func main() {

	wss := wss.NewWSServer("0.0.0.0", 9696, func(clientID string, messageType int, message []byte) {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("[%s] %s", red(clientID), string(message))
	})

	wss.Start()

	shell := NewShell(func(rl *readline.Instance, tokens []string) {
		command := tokens[0]
		args := tokens[1:]
		fn, ok := cmd.LookUp(command)
		if ok {
			fn(contextWith(rl, wss), args)
		}
	})

	err := shell.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

}

func contextWith(rl *readline.Instance, wss *wss.WSServer) context.Context {
	ctxWithWS := context.WithValue(context.Background(), constants.CtxKeyWss, wss)
	return context.WithValue(ctxWithWS, constants.CtxKeyReadline, rl)
}
