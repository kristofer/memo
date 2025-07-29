package cmd

import (
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