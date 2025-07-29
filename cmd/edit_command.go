package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"memo/internal/ui"
)

type EditCommand struct {
	ctx *CommandContext
}

func NewEditCommand(ctx *CommandContext) *EditCommand {
	return &EditCommand{ctx: ctx}
}

func (c *EditCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("note-id or number required\nUsage: memo edit <note-id|number>")
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

	fmt.Printf("Editing note: %s\n", n.Metadata.Title)
	fmt.Printf("Current content:\n%s\n\n", n.Content)

	newContent := ui.PromptForInput("Enter new content (leave empty to keep current): ")
	if newContent != "" {
		n.UpdateContent(newContent)
	}

	currentTags := strings.Join(n.Metadata.Tags, ", ")
	fmt.Printf("Current tags: %s\n", currentTags)
	newTags := ui.PromptForInput("Enter new tags (comma-separated, leave empty to keep current): ")
	if newTags != "" {
		var tags []string
		for _, tag := range strings.Split(newTags, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
		n.UpdateTags(tags)
	}

	err = c.ctx.Storage.SaveNote(n)
	if err != nil {
		return fmt.Errorf("error saving note: %w", err)
	}

	fmt.Println("Note updated successfully!")
	return nil
}

func (c *EditCommand) resolveNoteID(identifier string) (string, error) {
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