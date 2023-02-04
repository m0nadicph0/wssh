package cmd

import (
	"context"
	"github.com/chzyer/readline"
	"wssh/constants"
	"wssh/wss"
)

var handlerMap = map[string]Handler{
	"prompt": setPrompt,
	"send":   send,
	"list":   listClients,
}

type Handler func(ctx context.Context, args []string)

func LookUp(cmd string) (Handler, bool) {
	h, ok := handlerMap[cmd]
	return h, ok
}

func getInstance(ctx context.Context) *readline.Instance {
	return ctx.Value(constants.CtxKeyReadline).(*readline.Instance)
}

func getWS(ctx context.Context) *wss.WSServer {
	return ctx.Value(constants.CtxKeyWss).(*wss.WSServer)
}
