package cmd

import (
	"fmt"

	"memo/internal/ui"
)

type ReadCommand struct {
	ctx *CommandContext
}

func NewReadCommand(ctx *CommandContext) *ReadCommand {
	return &ReadCommand{ctx: ctx}
}

func (c *ReadCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("note-id or number required\nUsage: memo read <note-id|number>")
	}

	noteID, err := c.ctx.ResolveNoteID(args[0])
	if err != nil {
		return err
	}

	n, err := c.ctx.Storage.FindNoteByID(noteID)
	if err != nil {
		return err
	}

	ui.DisplayNote(n)
	return nil
}