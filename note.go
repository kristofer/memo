package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type NoteMetadata struct {
	Title    string    `yaml:"title"`
	Created  time.Time `yaml:"created"`
	Modified time.Time `yaml:"modified"`
	Tags     []string  `yaml:"tags,omitempty"`
	Author   string    `yaml:"author,omitempty"`
	Status   string    `yaml:"status,omitempty"`
	Priority int       `yaml:"priority,omitempty"`
}

type Note struct {
	Metadata NoteMetadata
	Content  string
	FilePath string
}

const (
	NotesDir      = ".memo-notes"
	NoteExtension = ".note"
)

func ensureNotesDir() error {
	if _, err := os.Stat(NotesDir); os.IsNotExist(err) {
		return os.MkdirAll(NotesDir, 0755)
	}
	return nil
}

func parseNote(filePath string) (*Note, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)
	
	// Check if file starts with YAML front matter
	if !strings.HasPrefix(contentStr, "---\n") {
		return nil, fmt.Errorf("note file must start with YAML front matter")
	}

	// Find the end of YAML front matter
	parts := strings.Split(contentStr, "\n---\n")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid note format: missing YAML front matter delimiter")
	}

	yamlContent := parts[0][4:] // Remove the first "---\n"
	noteContent := strings.Join(parts[1:], "\n---\n")

	var metadata NoteMetadata
	err = yaml.Unmarshal([]byte(yamlContent), &metadata)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML metadata: %v", err)
	}

	return &Note{
		Metadata: metadata,
		Content:  strings.TrimSpace(noteContent),
		FilePath: filePath,
	}, nil
}

func (n *Note) Save() error {
	if err := ensureNotesDir(); err != nil {
		return err
	}

	// Update modified time
	n.Metadata.Modified = time.Now()

	// Marshal YAML metadata
	yamlData, err := yaml.Marshal(&n.Metadata)
	if err != nil {
		return err
	}

	// Combine YAML header with content
	fullContent := fmt.Sprintf("---\n%s---\n\n%s", string(yamlData), n.Content)

	return os.WriteFile(n.FilePath, []byte(fullContent), 0644)
}

func generateNoteID() string {
	return fmt.Sprintf("note_%d", time.Now().Unix())
}

func generateNoteFilePath(noteID string) string {
	return filepath.Join(NotesDir, noteID+NoteExtension)
}

func getAllNotes() ([]*Note, error) {
	if err := ensureNotesDir(); err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(NotesDir, "*"+NoteExtension))
	if err != nil {
		return nil, err
	}

	var notes []*Note
	for _, file := range files {
		note, err := parseNote(file)
		if err != nil {
			fmt.Printf("Warning: failed to parse note %s: %v\n", file, err)
			continue
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func findNoteByID(noteID string) (*Note, error) {
	notePath := generateNoteFilePath(noteID)
	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("note with ID '%s' not found", noteID)
	}
	return parseNote(notePath)
}

func searchNotes(query string) ([]*Note, error) {
	notes, err := getAllNotes()
	if err != nil {
		return nil, err
	}

	var matches []*Note
	queryLower := strings.ToLower(query)
	
	for _, note := range notes {
		// Search in title, content, and tags
		if strings.Contains(strings.ToLower(note.Metadata.Title), queryLower) ||
			strings.Contains(strings.ToLower(note.Content), queryLower) {
			matches = append(matches, note)
			continue
		}
		
		// Search in tags
		for _, tag := range note.Metadata.Tags {
			if strings.Contains(strings.ToLower(tag), queryLower) {
				matches = append(matches, note)
				break
			}
		}
	}

	return matches, nil
}

func filterNotesByTag(tag string) ([]*Note, error) {
	notes, err := getAllNotes()
	if err != nil {
		return nil, err
	}

	var matches []*Note
	tagLower := strings.ToLower(tag)
	
	for _, note := range notes {
		for _, noteTag := range note.Metadata.Tags {
			if strings.ToLower(noteTag) == tagLower {
				matches = append(matches, note)
				break
			}
		}
	}

	return matches, nil
}

func openEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // Default editor
	}

	// For now, just return an error message
	fmt.Printf("Please edit the file manually: %s\n", filePath)
	fmt.Printf("You can set the EDITOR environment variable to use your preferred editor.\n")
	return nil
}

func promptForInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}