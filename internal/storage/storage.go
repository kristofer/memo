package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"memo/internal/note"
)

const (
	DefaultNotesDir      = ".memo-notes"
	DefaultNoteExtension = ".note"
)

type FileStorage struct {
	notesDir      string
	noteExtension string
}

func NewFileStorage() *FileStorage {
	return &FileStorage{
		notesDir:      DefaultNotesDir,
		noteExtension: DefaultNoteExtension,
	}
}

func NewFileStorageWithConfig(notesDir, noteExtension string) *FileStorage {
	return &FileStorage{
		notesDir:      notesDir,
		noteExtension: noteExtension,
	}
}

func (fs *FileStorage) EnsureNotesDir() error {
	if _, err := os.Stat(fs.notesDir); os.IsNotExist(err) {
		return os.MkdirAll(fs.notesDir, 0755)
	}
	return nil
}

func (fs *FileStorage) GenerateNoteID() string {
	return fmt.Sprintf("note_%d", time.Now().Unix())
}

func (fs *FileStorage) GenerateNoteFilePath(noteID string) string {
	return filepath.Join(fs.notesDir, noteID+fs.noteExtension)
}

func (fs *FileStorage) ParseNote(filePath string) (*note.Note, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	contentStr := string(content)

	if !strings.HasPrefix(contentStr, "---\n") {
		return nil, fmt.Errorf("note file must start with YAML front matter")
	}

	parts := strings.Split(contentStr, "\n---\n")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid note format: missing YAML front matter delimiter")
	}

	yamlContent := parts[0][4:] // Remove the first "---\n"
	noteContent := strings.Join(parts[1:], "\n---\n")

	var metadata note.Metadata
	err = yaml.Unmarshal([]byte(yamlContent), &metadata)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML metadata: %w", err)
	}

	n := &note.Note{
		Metadata: metadata,
		Content:  strings.TrimSpace(noteContent),
		FilePath: filePath,
	}

	return n, nil
}

func (fs *FileStorage) SaveNote(n *note.Note) error {
	if err := fs.EnsureNotesDir(); err != nil {
		return fmt.Errorf("error ensuring notes directory: %w", err)
	}

	return n.Save()
}

func (fs *FileStorage) GetAllNotes() ([]*note.Note, error) {
	if err := fs.EnsureNotesDir(); err != nil {
		return nil, fmt.Errorf("error ensuring notes directory: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(fs.notesDir, "*"+fs.noteExtension))
	if err != nil {
		return nil, fmt.Errorf("error finding note files: %w", err)
	}

	var notes []*note.Note
	for _, file := range files {
		n, err := fs.ParseNote(file)
		if err != nil {
			fmt.Printf("Warning: failed to parse note %s: %v\n", file, err)
			continue
		}
		notes = append(notes, n)
	}

	return notes, nil
}

func (fs *FileStorage) FindNoteByID(noteID string) (*note.Note, error) {
	notePath := fs.GenerateNoteFilePath(noteID)
	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("note with ID '%s' not found", noteID)
	}
	return fs.ParseNote(notePath)
}

func (fs *FileStorage) DeleteNote(noteID string) error {
	notePath := fs.GenerateNoteFilePath(noteID)
	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		return fmt.Errorf("note with ID '%s' not found", noteID)
	}
	return os.Remove(notePath)
}

func (fs *FileStorage) SearchNotes(query string) ([]*note.Note, error) {
	notes, err := fs.GetAllNotes()
	if err != nil {
		return nil, err
	}

	var matches []*note.Note
	queryLower := strings.ToLower(query)

	for _, n := range notes {
		if strings.Contains(strings.ToLower(n.Metadata.Title), queryLower) ||
			strings.Contains(strings.ToLower(n.Content), queryLower) {
			matches = append(matches, n)
			continue
		}

		for _, tag := range n.Metadata.Tags {
			if strings.Contains(strings.ToLower(tag), queryLower) {
				matches = append(matches, n)
				break
			}
		}
	}

	return matches, nil
}

func (fs *FileStorage) FilterNotesByTag(tag string) ([]*note.Note, error) {
	notes, err := fs.GetAllNotes()
	if err != nil {
		return nil, err
	}

	var matches []*note.Note
	tagLower := strings.ToLower(tag)

	for _, n := range notes {
		for _, noteTag := range n.Metadata.Tags {
			if strings.ToLower(noteTag) == tagLower {
				matches = append(matches, n)
				break
			}
		}
	}

	return matches, nil
}