package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

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

	identifier := args[0]
	noteID, err := c.resolveNoteID(identifier)
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

func (c *ReadCommand) resolveNoteID(identifier string) (string, error) {
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