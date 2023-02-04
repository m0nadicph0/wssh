package main

import (
	"github.com/chzyer/readline"
	"io"
	"log"
	"strings"
	"wssh/constants"
)

const TokenSeparator = " "

type LineHandler func(rl *readline.Instance, tokens []string)

type Shell struct {
	Handler LineHandler
}

func NewShell(fn LineHandler) *Shell {
	return &Shell{
		Handler: fn,
	}
}

func (sh Shell) Start() error {
	readlineInstance, err := readline.NewEx(&readline.Config{
		Prompt:          constants.Prompt,
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})

	if err != nil {
		return err
	}

	defer readlineInstance.Close()
	readlineInstance.CaptureExitSignal()

	log.SetOutput(readlineInstance.Stderr())
	for {
		line, err := readlineInstance.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		sh.Handler(readlineInstance, strings.Split(line, TokenSeparator))

	}

	return nil
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("mode",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
