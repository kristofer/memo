package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"memo/internal/note"
	"memo/internal/storage"
	"memo/internal/ui"
)

type ListCommand struct {
	ctx *CommandContext
}

func NewListCommand(ctx *CommandContext) *ListCommand {
	return &ListCommand{ctx: ctx}
}

func (c *ListCommand) Execute(args []string) error {
	var tagFilter string
	if len(args) >= 2 && args[0] == "--tag" {
		tagFilter = args[1]
	} else if len(args) >= 1 && args[0] == "--tag" {
		return fmt.Errorf("tag value required\nUsage: memo list --tag <tag>")
	}

	var notes []*note.Note
	var err error

	if tagFilter != "" {
		notes, err = c.ctx.Storage.FilterNotesByTag(tagFilter)
		if err != nil {
			return fmt.Errorf("error filtering notes by tag: %w", err)
		}
		fmt.Printf("Notes with tag '%s':\n", tagFilter)
	} else {
		notes, err = c.ctx.Storage.GetAllNotes()
		if err != nil {
			return fmt.Errorf("error listing notes: %w", err)
		}
		fmt.Println("All notes:")
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	// Update current listing for number-based access
	c.ctx.SetCurrentListing(notes)

	// Persist the listing so numbers remain valid in subsequent commands.
	noteIDs := make([]string, len(notes))
	for i, n := range notes {
		noteIDs[i] = strings.TrimSuffix(filepath.Base(n.FilePath), storage.DefaultNoteExtension)
	}
	if err := c.ctx.Storage.SaveLastListing(noteIDs); err != nil {
		fmt.Printf("Warning: failed to save listing for number-based access: %v\n", err)
	}

	ui.DisplayNotesWithPagination(notes)
	
	return nil
}