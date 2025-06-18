package main

import (
	"fmt"
	"os"
	"path/filepath" 
	"strconv"
	"strings"
	"time"
)

// Global variable to store current listing for number-based access
var currentListing []*Note

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "create":
		handleCreate()
	case "list":
		handleList()
	case "read":
		handleRead()
	case "edit":
		handleEdit()
	case "delete":
		handleDelete()
	case "search":
		handleSearch()
	case "stats":
		handleStats()
	case "--help", "-h", "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
	}
}

func printHelp() {
	fmt.Println("Memo - Personal Notes Manager")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  memo create                     Create a new note")
	fmt.Println("  memo list                       List all notes (with numbered references)")
	fmt.Println("  memo list --tag <tag>           List notes with specific tag")
	fmt.Println("  memo read <note-id|number>      Display a specific note")
	fmt.Println("  memo edit <note-id|number>      Edit a specific note") 
	fmt.Println("  memo delete <note-id|number>    Delete a specific note")
	fmt.Println("  memo search <query>             Search notes for text")
	fmt.Println("  memo stats                      Display statistics about your notes")
	fmt.Println("  memo --help                     Display this help information")
	fmt.Println("")
	fmt.Println("Note: After running 'memo list', you can use numbers 1-N to reference notes")
	fmt.Println("      instead of the full note ID (e.g., 'memo read 3' or 'memo edit 5')")
}

func handleCreate() {
	title := promptForInput("Enter note title: ")
	if title == "" {
		fmt.Println("Error: title is required")
		return
	}

	content := promptForInput("Enter note content: ")
	
	// Parse tags if provided
	tagsInput := promptForInput("Enter tags (comma-separated, optional): ")
	var tags []string
	if tagsInput != "" {
		for _, tag := range strings.Split(tagsInput, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	noteID := generateNoteID()
	note := &Note{
		Metadata: NoteMetadata{
			Title:    title,
			Created:  time.Now(),
			Modified: time.Now(),
			Tags:     tags,
		},
		Content:  content,
		FilePath: generateNoteFilePath(noteID),
	}

	err := note.Save()
	if err != nil {
		fmt.Printf("Error creating note: %v\n", err)
		return
	}

	fmt.Printf("Note created successfully: %s\n", noteID)
}

func handleList() {
	// Check for --tag flag
	var tagFilter string
	if len(os.Args) >= 3 && os.Args[2] == "--tag" {
		if len(os.Args) < 4 {
			fmt.Println("Error: tag value required")
			fmt.Println("Usage: memo list --tag <tag>")
			return
		}
		tagFilter = os.Args[3]
	}

	var notes []*Note
	var err error

	if tagFilter != "" {
		notes, err = filterNotesByTag(tagFilter)
		if err != nil {
			fmt.Printf("Error filtering notes by tag: %v\n", err)
			return
		}
		fmt.Printf("Notes with tag '%s':\n", tagFilter)
	} else {
		notes, err = getAllNotes()
		if err != nil {
			fmt.Printf("Error listing notes: %v\n", err)
			return
		}
		fmt.Println("All notes:")
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return
	}

	// Store notes in global listing for number-based access
	currentListing = notes
	
	// Display notes with pagination
	displayNotesWithPagination(notes)
}

func displayNotesWithPagination(notes []*Note) {
	const pageSize = 10
	startIndex := 0
	
	for {
		endIndex := startIndex + pageSize
		if endIndex > len(notes) {
			endIndex = len(notes)
		}
		
		// Display current page
		fmt.Printf("\nShowing notes %d-%d of %d:\n", startIndex+1, endIndex, len(notes))
		fmt.Println("========================================")
		
		for i := startIndex; i < endIndex; i++ {
			note := notes[i]
			noteID := strings.TrimSuffix(filepath.Base(note.FilePath), NoteExtension)
			listNumber := i + 1
			
			fmt.Printf("%2d. %s | Created: %s\n", 
				listNumber,
				note.Metadata.Title, 
				note.Metadata.Created.Format("2006-01-02 15:04"))
			
			if len(note.Metadata.Tags) > 0 {
				fmt.Printf("    Tags: %s\n", strings.Join(note.Metadata.Tags, ", "))
			}
			fmt.Printf("    ID: %s\n", noteID)
			fmt.Println()
		}
		
		// Check if there are more notes to display
		if endIndex >= len(notes) {
			fmt.Println("End of notes.")
			break
		}
		
		// Ask user if they want to see more
		fmt.Printf("Show next %d notes? (y/N): ", pageSize)
		response := promptForInput("")
		
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			break
		}
		
		startIndex = endIndex
	}
	
	fmt.Println("\nTip: Use 'memo read <number>' or 'memo edit <number>' with numbers 1-" + strconv.Itoa(len(notes)) + " from this listing.")
}

// resolveNoteID converts a note identifier (either actual ID or list number) to actual note ID
func resolveNoteID(identifier string) (string, error) {
	// Try to parse as a number first
	if num, err := strconv.Atoi(identifier); err == nil {
		// It's a number, check if it's valid for current listing
		if currentListing == nil || len(currentListing) == 0 {
			return "", fmt.Errorf("no current note listing. Please run 'memo list' first")
		}
		
		if num < 1 || num > len(currentListing) {
			return "", fmt.Errorf("number %d is out of range. Valid range: 1-%d", num, len(currentListing))
		}
		
		// Convert to actual note ID
		note := currentListing[num-1]
		return strings.TrimSuffix(filepath.Base(note.FilePath), NoteExtension), nil
	}
	
	// It's not a number, assume it's an actual note ID
	return identifier, nil
}

func handleRead() {
	if len(os.Args) < 3 {
		fmt.Println("Error: note-id or number required")
		fmt.Println("Usage: memo read <note-id|number>")
		return
	}
	identifier := os.Args[2]
	
	// Resolve the identifier to actual note ID
	noteID, err := resolveNoteID(identifier)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	note, err := findNoteByID(noteID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Title: %s\n", note.Metadata.Title)
	fmt.Printf("Created: %s\n", note.Metadata.Created.Format("2006-01-02 15:04:05"))
	fmt.Printf("Modified: %s\n", note.Metadata.Modified.Format("2006-01-02 15:04:05"))
	
	if len(note.Metadata.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(note.Metadata.Tags, ", "))
	}
	
	if note.Metadata.Author != "" {
		fmt.Printf("Author: %s\n", note.Metadata.Author)
	}
	
	if note.Metadata.Status != "" {
		fmt.Printf("Status: %s\n", note.Metadata.Status)
	}
	
	if note.Metadata.Priority > 0 {
		fmt.Printf("Priority: %d\n", note.Metadata.Priority)
	}
	
	fmt.Println("\nContent:")
	fmt.Println("--------")
	fmt.Println(note.Content)
}

func handleEdit() {
	if len(os.Args) < 3 {
		fmt.Println("Error: note-id or number required")
		fmt.Println("Usage: memo edit <note-id|number>")
		return
	}
	identifier := os.Args[2]
	
	// Resolve the identifier to actual note ID
	noteID, err := resolveNoteID(identifier)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	note, err := findNoteByID(noteID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Editing note: %s\n", note.Metadata.Title)
	fmt.Printf("Current content:\n%s\n\n", note.Content)
	
	// Simple edit - just replace content
	newContent := promptForInput("Enter new content (leave empty to keep current): ")
	if newContent != "" {
		note.Content = newContent
	}
	
	// Update tags
	currentTags := strings.Join(note.Metadata.Tags, ", ")
	fmt.Printf("Current tags: %s\n", currentTags)
	newTags := promptForInput("Enter new tags (comma-separated, leave empty to keep current): ")
	if newTags != "" {
		var tags []string
		for _, tag := range strings.Split(newTags, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
		note.Metadata.Tags = tags
	}

	err = note.Save()
	if err != nil {
		fmt.Printf("Error saving note: %v\n", err)
		return
	}

	fmt.Println("Note updated successfully!")
}

func handleDelete() {
	if len(os.Args) < 3 {
		fmt.Println("Error: note-id or number required")
		fmt.Println("Usage: memo delete <note-id|number>")
		return
	}
	identifier := os.Args[2]
	
	// Resolve the identifier to actual note ID
	noteID, err := resolveNoteID(identifier)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	note, err := findNoteByID(noteID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Are you sure you want to delete note '%s'? (y/N): ", note.Metadata.Title)
	confirmation := promptForInput("")
	
	if strings.ToLower(confirmation) != "y" && strings.ToLower(confirmation) != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}

	err = os.Remove(note.FilePath)
	if err != nil {
		fmt.Printf("Error deleting note: %v\n", err)
		return
	}

	fmt.Println("Note deleted successfully!")
}

func handleSearch() {
	if len(os.Args) < 3 {
		fmt.Println("Error: search query required")
		fmt.Println("Usage: memo search <query>")
		return
	}
	query := os.Args[2]
	
	notes, err := searchNotes(query)
	if err != nil {
		fmt.Printf("Error searching notes: %v\n", err)
		return
	}

	if len(notes) == 0 {
		fmt.Printf("No notes found matching '%s'\n", query)
		return
	}

	fmt.Printf("Found %d note(s) matching '%s':\n\n", len(notes), query)
	
	for _, note := range notes {
		noteID := strings.TrimSuffix(filepath.Base(note.FilePath), NoteExtension)
		fmt.Printf("ID: %s | Title: %s\n", noteID, note.Metadata.Title)
		
		// Show preview of content
		preview := note.Content
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		fmt.Printf("Preview: %s\n", preview)
		fmt.Println("--------")
	}
}

func handleStats() {
	notes, err := getAllNotes()
	if err != nil {
		fmt.Printf("Error loading notes: %v\n", err)
		return
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return
	}

	fmt.Println("Note Statistics:")
	fmt.Printf("Total notes: %d\n", len(notes))
	
	// Count tags
	tagCount := make(map[string]int)
	var totalWords int
	var oldestNote, newestNote *Note
	
	for i, note := range notes {
		// Count words in content
		words := strings.Fields(note.Content)
		totalWords += len(words)
		
		// Track oldest and newest notes
		if i == 0 {
			oldestNote = note
			newestNote = note
		} else {
			if note.Metadata.Created.Before(oldestNote.Metadata.Created) {
				oldestNote = note
			}
			if note.Metadata.Created.After(newestNote.Metadata.Created) {
				newestNote = note
			}
		}
		
		// Count tags
		for _, tag := range note.Metadata.Tags {
			tagCount[tag]++
		}
	}
	
	fmt.Printf("Total words: %d\n", totalWords)
	fmt.Printf("Average words per note: %.1f\n", float64(totalWords)/float64(len(notes)))
	
	if oldestNote != nil {
		fmt.Printf("Oldest note: %s (%s)\n", oldestNote.Metadata.Title, oldestNote.Metadata.Created.Format("2006-01-02"))
	}
	if newestNote != nil {
		fmt.Printf("Newest note: %s (%s)\n", newestNote.Metadata.Title, newestNote.Metadata.Created.Format("2006-01-02"))
	}
	
	if len(tagCount) > 0 {
		fmt.Printf("\nTag usage:\n")
		for tag, count := range tagCount {
			fmt.Printf("  %s: %d\n", tag, count)
		}
	}
}