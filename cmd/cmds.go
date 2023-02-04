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
	for _, clientID := range getWS(ctx).GetClientIDs() {
		color.Yellow(clientID)
	}
}

func broadcast(ctx context.Context, args []string) {
	data := strings.Join(args, " ")
	getWS(ctx).BroadcastText(data)
}
