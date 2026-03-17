package cmd

import (
	"fmt"

	"memo/internal/ui"
)

type DeleteCommand struct {
	ctx *CommandContext
}

func NewDeleteCommand(ctx *CommandContext) *DeleteCommand {
	return &DeleteCommand{ctx: ctx}
}

func (c *DeleteCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("note-id or number required\nUsage: memo delete <note-id|number>")
	}

	noteID, err := c.ctx.ResolveNoteID(args[0])
	if err != nil {
		return err
	}

	n, err := c.ctx.Storage.FindNoteByID(noteID)
	if err != nil {
		return err
	}

	prompt := fmt.Sprintf("Are you sure you want to delete note '%s'? (y/N): ", n.Metadata.Title)
	if !ui.ConfirmAction(prompt) {
		fmt.Println("Deletion cancelled.")
		return nil
	}

	err = c.ctx.Storage.DeleteNote(noteID)
	if err != nil {
		return fmt.Errorf("error deleting note: %w", err)
	}

	fmt.Println("Note deleted successfully!")
	return nil
}