package note

import (
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	title := "Test Note"
	content := "Test content"
	tags := []string{"go", "testing"}

	n := New(title, content, tags)

	if n.Metadata.Title != title {
		t.Errorf("expected title %q, got %q", title, n.Metadata.Title)
	}
	if n.Content != content {
		t.Errorf("expected content %q, got %q", content, n.Content)
	}
	if len(n.Metadata.Tags) != 2 || n.Metadata.Tags[0] != "go" || n.Metadata.Tags[1] != "testing" {
		t.Errorf("unexpected tags: %v", n.Metadata.Tags)
	}
	if n.Metadata.Created.IsZero() {
		t.Error("expected non-zero Created time")
	}
	if n.Metadata.Modified.IsZero() {
		t.Error("expected non-zero Modified time")
	}
}

func TestNew_NoTags(t *testing.T) {
	n := New("title", "body", nil)
	if len(n.Metadata.Tags) != 0 {
		t.Errorf("expected no tags, got %v", n.Metadata.Tags)
	}
}

func TestUpdateContent(t *testing.T) {
	n := New("title", "old content", nil)
	before := n.Metadata.Modified

	// Ensure measurable time difference
	time.Sleep(time.Millisecond)
	n.UpdateContent("new content")

	if n.Content != "new content" {
		t.Errorf("expected content %q, got %q", "new content", n.Content)
	}
	if !n.Metadata.Modified.After(before) {
		t.Error("expected Modified to be updated after UpdateContent")
	}
}

func TestUpdateTags(t *testing.T) {
	n := New("title", "content", []string{"old"})
	time.Sleep(time.Millisecond)
	n.UpdateTags([]string{"new1", "new2"})

	if len(n.Metadata.Tags) != 2 || n.Metadata.Tags[0] != "new1" {
		t.Errorf("unexpected tags: %v", n.Metadata.Tags)
	}
}

func TestToFileContent(t *testing.T) {
	n := New("My Note", "Hello world", []string{"a", "b"})

	content, err := n.ToFileContent()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(content, "---\n") {
		t.Error("file content should start with YAML front matter delimiter")
	}
	if !strings.Contains(content, "title: My Note") {
		t.Error("file content should contain title in YAML front matter")
	}
	if !strings.Contains(content, "Hello world") {
		t.Error("file content should contain note body")
	}
}

func TestSetFilePath(t *testing.T) {
	n := New("title", "content", nil)
	n.SetFilePath("/tmp/test.note")
	if n.FilePath != "/tmp/test.note" {
		t.Errorf("expected FilePath %q, got %q", "/tmp/test.note", n.FilePath)
	}
}
