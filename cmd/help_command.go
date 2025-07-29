package cmd

import "memo/internal/ui"

type HelpCommand struct {
	ctx *CommandContext
}

func NewHelpCommand(ctx *CommandContext) *HelpCommand {
	return &HelpCommand{ctx: ctx}
}

func (c *HelpCommand) Execute(args []string) error {
	ui.PrintHelp()
	return nil
}