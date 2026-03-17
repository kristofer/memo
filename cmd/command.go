package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"memo/internal/note"
	"memo/internal/storage"
)

// Command interface defines the contract for all CLI commands
type Command interface {
	Execute(args []string) error
}

// CommandContext provides shared dependencies for all commands
type CommandContext struct {
	Storage        *storage.FileStorage
	CurrentListing []*note.Note
}

// SetCurrentListing updates the current listing (used by list command)
func (ctx *CommandContext) SetCurrentListing(notes []*note.Note) {
	ctx.CurrentListing = notes
}

// GetCurrentListing returns the current listing
func (ctx *CommandContext) GetCurrentListing() []*note.Note {
	return ctx.CurrentListing
}

// ResolveNoteID converts a numeric list position or a literal note ID string
// into a canonical note ID (the filename without its extension).
// Numeric positions require a prior call to the list command to populate CurrentListing.
func (ctx *CommandContext) ResolveNoteID(identifier string) (string, error) {
	if num, err := strconv.Atoi(identifier); err == nil {
		if len(ctx.CurrentListing) == 0 {
			return "", fmt.Errorf("no current note listing. Please run 'memo list' first")
		}
		if num < 1 || num > len(ctx.CurrentListing) {
			return "", fmt.Errorf("number %d is out of range. Valid range: 1-%d", num, len(ctx.CurrentListing))
		}
		n := ctx.CurrentListing[num-1]
		return strings.TrimSuffix(filepath.Base(n.FilePath), storage.DefaultNoteExtension), nil
	}
	return identifier, nil
}