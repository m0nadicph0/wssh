package cmd

import (
	"context"
	"github.com/fatih/color"
	"strings"
	"wssh/constants"
)

func setPrompt(ctx context.Context, args []string) {
	rl := getInstance(ctx)
	rl.Refresh()
	switch len(args) {
	case 0:
		rl.SetPrompt(constants.Prompt)
	default:
		rl.SetPrompt(strings.Join(args, " ") + " ")
	}
}

func send(ctx context.Context, args []string) {
	clientID := args[0]
	data := strings.Join(args[1:], " ")
	getWS(ctx).WriteText(clientID, data)
}

func listClients(ctx context.Context, args []string) {
	wss := getWS(ctx)
	cids := wss.GetClientIDs()
	for _, clientID := range cids {
		color.Blue(clientID)
	}
}
