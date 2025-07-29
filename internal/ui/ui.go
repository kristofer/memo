package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"memo/internal/note"
)

func PromptForInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func PrintHelp() {
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

func DisplayNotesWithPagination(notes []*note.Note) {
	const pageSize = 10
	startIndex := 0

	for {
		endIndex := startIndex + pageSize
		if endIndex > len(notes) {
			endIndex = len(notes)
		}

		fmt.Printf("\nShowing notes %d-%d of %d:\n", startIndex+1, endIndex, len(notes))
		fmt.Println("========================================")

		for i := startIndex; i < endIndex; i++ {
			n := notes[i]
			noteID := strings.TrimSuffix(filepath.Base(n.FilePath), ".note")
			listNumber := i + 1

			fmt.Printf("%2d. %s | Created: %s\n",
				listNumber,
				n.Metadata.Title,
				n.Metadata.Created.Format("2006-01-02 15:04"))

			if len(n.Metadata.Tags) > 0 {
				fmt.Printf("    Tags: %s\n", strings.Join(n.Metadata.Tags, ", "))
			}
			fmt.Printf("    ID: %s\n", noteID)
			fmt.Println()
		}

		if endIndex >= len(notes) {
			fmt.Println("End of notes.")
			break
		}

		fmt.Printf("Show next %d notes? (y/N): ", pageSize)
		response := PromptForInput("")

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			break
		}

		startIndex = endIndex
	}

	fmt.Println("\nTip: Use 'memo read <number>' or 'memo edit <number>' with numbers 1-" + strconv.Itoa(len(notes)) + " from this listing.")
}

func DisplayNote(n *note.Note) {
	fmt.Printf("Title: %s\n", n.Metadata.Title)
	fmt.Printf("Created: %s\n", n.Metadata.Created.Format("2006-01-02 15:04:05"))
	fmt.Printf("Modified: %s\n", n.Metadata.Modified.Format("2006-01-02 15:04:05"))

	if len(n.Metadata.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(n.Metadata.Tags, ", "))
	}

	if n.Metadata.Author != "" {
		fmt.Printf("Author: %s\n", n.Metadata.Author)
	}

	if n.Metadata.Status != "" {
		fmt.Printf("Status: %s\n", n.Metadata.Status)
	}

	if n.Metadata.Priority > 0 {
		fmt.Printf("Priority: %d\n", n.Metadata.Priority)
	}

	fmt.Println("\nContent:")
	fmt.Println("--------")
	fmt.Println(n.Content)
}

func DisplaySearchResults(notes []*note.Note, query string) {
	if len(notes) == 0 {
		fmt.Printf("No notes found matching '%s'\n", query)
		return
	}

	fmt.Printf("Found %d note(s) matching '%s':\n\n", len(notes), query)

	for _, n := range notes {
		noteID := strings.TrimSuffix(filepath.Base(n.FilePath), ".note")
		fmt.Printf("ID: %s | Title: %s\n", noteID, n.Metadata.Title)

		preview := n.Content
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		fmt.Printf("Preview: %s\n", preview)
		fmt.Println("--------")
	}
}

func DisplayStats(notes []*note.Note) {
	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return
	}

	fmt.Println("Note Statistics:")
	fmt.Printf("Total notes: %d\n", len(notes))

	tagCount := make(map[string]int)
	var totalWords int
	var oldestNote, newestNote *note.Note

	for i, n := range notes {
		words := strings.Fields(n.Content)
		totalWords += len(words)

		if i == 0 {
			oldestNote = n
			newestNote = n
		} else {
			if n.Metadata.Created.Before(oldestNote.Metadata.Created) {
				oldestNote = n
			}
			if n.Metadata.Created.After(newestNote.Metadata.Created) {
				newestNote = n
			}
		}

		for _, tag := range n.Metadata.Tags {
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

func ConfirmAction(prompt string) bool {
	response := PromptForInput(prompt)
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}