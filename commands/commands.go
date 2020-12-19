package commands

import (
	"fmt"
	"strings"

	"github.com/pavlovskyive/architecture-lab-4/engine"
)

type printCommand struct {
	arg string
}

func (pCmd *printCommand) Execute(_ engine.Handler) {
	fmt.Println(pCmd.arg)
}

type reverseCommand struct {
	arg string
}

func (rCmd *reverseCommand) Execute(h engine.Handler) {
	res := ""
	for _, a := range rCmd.arg {
		res = string(a) + res
	}
	h.Post(&printCommand{arg: res})
}

// Parse translates command lines into commands
func Parse(commandLine string) engine.Command {

	parts := strings.Fields(commandLine)

	if len(parts) != 2 {
		return &engine.printCommand{arg: "SYNTAX ERROR: wrong number of input elements in line: " + commandLine}
	}

	command := parts[0]
	arg := parts[1]

	switch command {
	case "print":
		return &printCommand{arg: arg}
	case "reverse":
		return &reverseCommand{arg: arg}
	default:
		return &printCommand{arg: "SYNTAX ERROR: unrecognizable command: " + command}
	}
}
