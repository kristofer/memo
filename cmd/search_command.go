package cmd

import (
	"fmt"

	"memo/internal/ui"
)

type SearchCommand struct {
	ctx *CommandContext
}

func NewSearchCommand(ctx *CommandContext) *SearchCommand {
	return &SearchCommand{ctx: ctx}
}

func (c *SearchCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("search query required\nUsage: memo search <query>")
	}

	query := args[0]
	notes, err := c.ctx.Storage.SearchNotes(query)
	if err != nil {
		return fmt.Errorf("error searching notes: %w", err)
	}

	ui.DisplaySearchResults(notes, query)
	return nil
}