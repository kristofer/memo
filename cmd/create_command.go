package cmd

import (
	"fmt"
	"strings"

	"memo/internal/note"
	"memo/internal/ui"
)

type CreateCommand struct {
	ctx *CommandContext
}

func NewCreateCommand(ctx *CommandContext) *CreateCommand {
	return &CreateCommand{ctx: ctx}
}

func (c *CreateCommand) Execute(args []string) error {
	title := ui.PromptForInput("Enter note title: ")
	if title == "" {
		return fmt.Errorf("title is required")
	}

	content := ui.PromptForInput("Enter note content: ")

	tagsInput := ui.PromptForInput("Enter tags (comma-separated, optional): ")
	var tags []string
	if tagsInput != "" {
		for _, tag := range strings.Split(tagsInput, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	noteID := c.ctx.Storage.GenerateNoteID()
	n := note.New(title, content, tags)
	n.SetFilePath(c.ctx.Storage.GenerateNoteFilePath(noteID))

	err := c.ctx.Storage.SaveNote(n)
	if err != nil {
		return fmt.Errorf("error creating note: %w", err)
	}

	fmt.Printf("Note created successfully: %s\n", noteID)
	return nil
}