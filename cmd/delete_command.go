package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

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

	identifier := args[0]
	noteID, err := c.resolveNoteID(identifier)
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

func (c *DeleteCommand) resolveNoteID(identifier string) (string, error) {
	if num, err := strconv.Atoi(identifier); err == nil {
		if c.ctx.CurrentListing == nil || len(c.ctx.CurrentListing) == 0 {
			return "", fmt.Errorf("no current note listing. Please run 'memo list' first")
		}

		if num < 1 || num > len(c.ctx.CurrentListing) {
			return "", fmt.Errorf("number %d is out of range. Valid range: 1-%d", num, len(c.ctx.CurrentListing))
		}

		n := c.ctx.CurrentListing[num-1]
		return strings.TrimSuffix(filepath.Base(n.FilePath), ".note"), nil
	}

	return identifier, nil
}