package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"memo/internal/note"
	"memo/internal/storage"
	"memo/internal/ui"
)

var currentListing []*note.Note

type App struct {
	storage *storage.FileStorage
}

func NewApp() *App {
	return &App{
		storage: storage.NewFileStorage(),
	}
}

func (app *App) HandleCreate() {
	title := ui.PromptForInput("Enter note title: ")
	if title == "" {
		fmt.Println("Error: title is required")
		return
	}

	content := ui.PromptForInput("Enter note content: ")

	tagsInput := ui.PromptForInput("Enter tags (comma-separated, optional): ")
	var tags []string
	if tagsInput != "" {
		for _, tag := range strings.Split(tagsInput, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	noteID := app.storage.GenerateNoteID()
	n := note.New(title, content, tags)
	n.SetFilePath(app.storage.GenerateNoteFilePath(noteID))

	err := app.storage.SaveNote(n)
	if err != nil {
		fmt.Printf("Error creating note: %v\n", err)
		return
	}

	fmt.Printf("Note created successfully: %s\n", noteID)
}

func (app *App) HandleList(args []string) {
	var tagFilter string
	if len(args) >= 2 && args[0] == "--tag" {
		if len(args) < 2 {
			fmt.Println("Error: tag value required")
			fmt.Println("Usage: memo list --tag <tag>")
			return
		}
		tagFilter = args[1]
	}

	var notes []*note.Note
	var err error

	if tagFilter != "" {
		notes, err = app.storage.FilterNotesByTag(tagFilter)
		if err != nil {
			fmt.Printf("Error filtering notes by tag: %v\n", err)
			return
		}
		fmt.Printf("Notes with tag '%s':\n", tagFilter)
	} else {
		notes, err = app.storage.GetAllNotes()
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

	currentListing = notes
	ui.DisplayNotesWithPagination(notes)
}

func (app *App) resolveNoteID(identifier string) (string, error) {
	if num, err := strconv.Atoi(identifier); err == nil {
		if currentListing == nil || len(currentListing) == 0 {
			return "", fmt.Errorf("no current note listing. Please run 'memo list' first")
		}

		if num < 1 || num > len(currentListing) {
			return "", fmt.Errorf("number %d is out of range. Valid range: 1-%d", num, len(currentListing))
		}

		n := currentListing[num-1]
		return strings.TrimSuffix(filepath.Base(n.FilePath), ".note"), nil
	}

	return identifier, nil
}

func (app *App) HandleRead(identifier string) {
	noteID, err := app.resolveNoteID(identifier)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	n, err := app.storage.FindNoteByID(noteID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ui.DisplayNote(n)
}

func (app *App) HandleEdit(identifier string) {
	noteID, err := app.resolveNoteID(identifier)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	n, err := app.storage.FindNoteByID(noteID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
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

	err = app.storage.SaveNote(n)
	if err != nil {
		fmt.Printf("Error saving note: %v\n", err)
		return
	}

	fmt.Println("Note updated successfully!")
}

func (app *App) HandleDelete(identifier string) {
	noteID, err := app.resolveNoteID(identifier)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	n, err := app.storage.FindNoteByID(noteID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	prompt := fmt.Sprintf("Are you sure you want to delete note '%s'? (y/N): ", n.Metadata.Title)
	if !ui.ConfirmAction(prompt) {
		fmt.Println("Deletion cancelled.")
		return
	}

	err = app.storage.DeleteNote(noteID)
	if err != nil {
		fmt.Printf("Error deleting note: %v\n", err)
		return
	}

	fmt.Println("Note deleted successfully!")
}

func (app *App) HandleSearch(query string) {
	notes, err := app.storage.SearchNotes(query)
	if err != nil {
		fmt.Printf("Error searching notes: %v\n", err)
		return
	}

	ui.DisplaySearchResults(notes, query)
}

func (app *App) HandleStats() {
	notes, err := app.storage.GetAllNotes()
	if err != nil {
		fmt.Printf("Error loading notes: %v\n", err)
		return
	}

	ui.DisplayStats(notes)
}

func (app *App) Run() {
	if len(os.Args) < 2 {
		ui.PrintHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "create":
		app.HandleCreate()
	case "list":
		var args []string
		if len(os.Args) > 2 {
			args = os.Args[2:]
		}
		app.HandleList(args)
	case "read":
		if len(os.Args) < 3 {
			fmt.Println("Error: note-id or number required")
			fmt.Println("Usage: memo read <note-id|number>")
			return
		}
		app.HandleRead(os.Args[2])
	case "edit":
		if len(os.Args) < 3 {
			fmt.Println("Error: note-id or number required")
			fmt.Println("Usage: memo edit <note-id|number>")
			return
		}
		app.HandleEdit(os.Args[2])
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: note-id or number required")
			fmt.Println("Usage: memo delete <note-id|number>")
			return
		}
		app.HandleDelete(os.Args[2])
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Error: search query required")
			fmt.Println("Usage: memo search <query>")
			return
		}
		app.HandleSearch(os.Args[2])
	case "stats":
		app.HandleStats()
	case "--help", "-h", "help":
		ui.PrintHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		ui.PrintHelp()
	}
}