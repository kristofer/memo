package cmd

import (
	"fmt"

	"memo/internal/ui"
)

type StatsCommand struct {
	ctx *CommandContext
}

func NewStatsCommand(ctx *CommandContext) *StatsCommand {
	return &StatsCommand{ctx: ctx}
}

func (c *StatsCommand) Execute(args []string) error {
	notes, err := c.ctx.Storage.GetAllNotes()
	if err != nil {
		return fmt.Errorf("error loading notes: %w", err)
	}

	ui.DisplayStats(notes)
	return nil
}