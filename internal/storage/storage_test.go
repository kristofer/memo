package storage

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"memo/internal/note"
)

// newTestStorage creates a FileStorage backed by a temporary directory.
func newTestStorage(t *testing.T) *FileStorage {
	t.Helper()
	dir := t.TempDir()
	return NewFileStorageWithConfig(dir, DefaultNoteExtension)
}

func TestGenerateNoteID_Format(t *testing.T) {
	fs := newTestStorage(t)
	id := fs.GenerateNoteID()

	// Expected format: YYYYMMDD-HHMMSS-mmm (19 characters)
	if len(id) != 19 {
		t.Errorf("expected note ID length 19, got %d (%q)", len(id), id)
	}
	if id[8] != '-' {
		t.Errorf("expected '-' at position 8, got %q in %q", string(id[8]), id)
	}
	if id[15] != '-' {
		t.Errorf("expected '-' at position 15, got %q in %q", string(id[15]), id)
	}
}

func TestGenerateNoteFilePath(t *testing.T) {
	fs := newTestStorage(t)
	path := fs.GenerateNoteFilePath("20260317-193750-042")
	if !strings.HasSuffix(path, "20260317-193750-042.note") {
		t.Errorf("unexpected note file path: %q", path)
	}
}

func TestEnsureNotesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "sub", "notes")
	fs := NewFileStorageWithConfig(dir, DefaultNoteExtension)

	if err := fs.EnsureNotesDir(); err != nil {
		t.Fatalf("EnsureNotesDir returned error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected notes directory to be created")
	}
}

func TestSaveAndParseNote(t *testing.T) {
	fs := newTestStorage(t)

	n := note.New("Save Test", "some content", []string{"tag1"})
	id := fs.GenerateNoteID()
	n.SetFilePath(fs.GenerateNoteFilePath(id))

	if err := fs.SaveNote(n); err != nil {
		t.Fatalf("SaveNote returned error: %v", err)
	}

	loaded, err := fs.ParseNote(n.FilePath)
	if err != nil {
		t.Fatalf("ParseNote returned error: %v", err)
	}

	if loaded.Metadata.Title != "Save Test" {
		t.Errorf("expected title %q, got %q", "Save Test", loaded.Metadata.Title)
	}
	if loaded.Content != "some content" {
		t.Errorf("expected content %q, got %q", "some content", loaded.Content)
	}
	if len(loaded.Metadata.Tags) != 1 || loaded.Metadata.Tags[0] != "tag1" {
		t.Errorf("unexpected tags: %v", loaded.Metadata.Tags)
	}
}

func TestParseNote_InvalidFormat(t *testing.T) {
	fs := newTestStorage(t)

	path := filepath.Join(fs.notesDir, "bad.note")
	if err := os.WriteFile(path, []byte("no front matter here"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := fs.ParseNote(path); err == nil {
		t.Error("expected error for note without YAML front matter")
	}
}

func TestGetAllNotes(t *testing.T) {
	fs := newTestStorage(t)

	for i, title := range []string{"Alpha", "Beta", "Gamma"} {
		n := note.New(title, "body", nil)
		id := fs.GenerateNoteID() + "-" + strconv.Itoa(i) // ensure uniqueness
		n.SetFilePath(fs.GenerateNoteFilePath(id))
		if err := fs.SaveNote(n); err != nil {
			t.Fatalf("SaveNote(%q): %v", title, err)
		}
	}

	notes, err := fs.GetAllNotes()
	if err != nil {
		t.Fatalf("GetAllNotes returned error: %v", err)
	}
	if len(notes) != 3 {
		t.Errorf("expected 3 notes, got %d", len(notes))
	}
}

func TestFindNoteByID(t *testing.T) {
	fs := newTestStorage(t)

	n := note.New("Find Me", "content", nil)
	id := fs.GenerateNoteID()
	n.SetFilePath(fs.GenerateNoteFilePath(id))
	if err := fs.SaveNote(n); err != nil {
		t.Fatal(err)
	}

	found, err := fs.FindNoteByID(id)
	if err != nil {
		t.Fatalf("FindNoteByID returned error: %v", err)
	}
	if found.Metadata.Title != "Find Me" {
		t.Errorf("expected title %q, got %q", "Find Me", found.Metadata.Title)
	}
}

func TestFindNoteByID_NotFound(t *testing.T) {
	fs := newTestStorage(t)
	if _, err := fs.FindNoteByID("nonexistent-id"); err == nil {
		t.Error("expected error for non-existent note ID")
	}
}

func TestDeleteNote(t *testing.T) {
	fs := newTestStorage(t)

	n := note.New("Delete Me", "content", nil)
	id := fs.GenerateNoteID()
	n.SetFilePath(fs.GenerateNoteFilePath(id))
	if err := fs.SaveNote(n); err != nil {
		t.Fatal(err)
	}

	if err := fs.DeleteNote(id); err != nil {
		t.Fatalf("DeleteNote returned error: %v", err)
	}

	if _, err := os.Stat(n.FilePath); !os.IsNotExist(err) {
		t.Error("expected note file to be removed after deletion")
	}
}

func TestDeleteNote_NotFound(t *testing.T) {
	fs := newTestStorage(t)
	if err := fs.DeleteNote("nonexistent-id"); err == nil {
		t.Error("expected error when deleting non-existent note")
	}
}

func TestSearchNotes(t *testing.T) {
	fs := newTestStorage(t)

	notes := []struct{ title, content string }{
		{"Go Programming", "goroutines and channels"},
		{"Python Basics", "lists and dicts"},
		{"Go Concurrency", "mutexes and waitgroups"},
	}
	for i, nd := range notes {
		n := note.New(nd.title, nd.content, nil)
		id := fs.GenerateNoteID() + "-" + strconv.Itoa(i) // ensure uniqueness
		n.SetFilePath(fs.GenerateNoteFilePath(id))
		if err := fs.SaveNote(n); err != nil {
			t.Fatal(err)
		}
	}

	results, err := fs.SearchNotes("go")
	if err != nil {
		t.Fatalf("SearchNotes returned error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results for 'go', got %d", len(results))
	}
}

func TestFilterNotesByTag(t *testing.T) {
	fs := newTestStorage(t)

	for i, tags := range [][]string{{"go", "backend"}, {"python"}, {"go", "cli"}} {
		n := note.New("Note", "body", tags)
		id := fs.GenerateNoteID() + "-" + strconv.Itoa(i) // ensure uniqueness
		n.SetFilePath(fs.GenerateNoteFilePath(id))
		if err := fs.SaveNote(n); err != nil {
			t.Fatal(err)
		}
	}

	results, err := fs.FilterNotesByTag("go")
	if err != nil {
		t.Fatalf("FilterNotesByTag returned error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 notes with tag 'go', got %d", len(results))
	}
}

func TestSaveAndLoadLastListing(t *testing.T) {
	fs := newTestStorage(t)

	ids := []string{"20260317-193750-001", "20260317-193750-002", "20260317-193750-003"}
	if err := fs.SaveLastListing(ids); err != nil {
		t.Fatalf("SaveLastListing returned error: %v", err)
	}

	loaded, err := fs.LoadLastListing()
	if err != nil {
		t.Fatalf("LoadLastListing returned error: %v", err)
	}
	if len(loaded) != len(ids) {
		t.Fatalf("expected %d IDs, got %d", len(ids), len(loaded))
	}
	for i, id := range ids {
		if loaded[i] != id {
			t.Errorf("position %d: expected %q, got %q", i+1, id, loaded[i])
		}
	}
}

func TestLoadLastListing_NoFile(t *testing.T) {
	fs := newTestStorage(t)

	ids, err := fs.LoadLastListing()
	if err != nil {
		t.Fatalf("LoadLastListing should not error when no listing file exists, got: %v", err)
	}
	if ids != nil {
		t.Errorf("expected nil when no listing file exists, got %v", ids)
	}
}

func TestSaveLastListing_Overwrites(t *testing.T) {
	fs := newTestStorage(t)

	first := []string{"id-a", "id-b", "id-c"}
	if err := fs.SaveLastListing(first); err != nil {
		t.Fatal(err)
	}

	second := []string{"id-x", "id-y"}
	if err := fs.SaveLastListing(second); err != nil {
		t.Fatal(err)
	}

	loaded, err := fs.LoadLastListing()
	if err != nil {
		t.Fatalf("LoadLastListing returned error: %v", err)
	}
	if len(loaded) != len(second) {
		t.Fatalf("expected %d IDs after overwrite, got %d", len(second), len(loaded))
	}
	for i, id := range second {
		if loaded[i] != id {
			t.Errorf("position %d: expected %q, got %q", i+1, id, loaded[i])
		}
	}
}

func TestDefaultNotesDir(t *testing.T) {
	dir := DefaultNotesDir()
	if dir == "" {
		t.Error("DefaultNotesDir should return a non-empty path")
	}
	if !strings.HasSuffix(dir, DefaultNotesDirName) {
		t.Errorf("expected path to end with %q, got %q", DefaultNotesDirName, dir)
	}
}
