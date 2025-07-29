package note

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Title    string    `yaml:"title"`
	Created  time.Time `yaml:"created"`
	Modified time.Time `yaml:"modified"`
	Tags     []string  `yaml:"tags,omitempty"`
	Author   string    `yaml:"author,omitempty"`
	Status   string    `yaml:"status,omitempty"`
	Priority int       `yaml:"priority,omitempty"`
}

type Note struct {
	Metadata Metadata
	Content  string
	FilePath string
}

func New(title, content string, tags []string) *Note {
	now := time.Now()
	return &Note{
		Metadata: Metadata{
			Title:    title,
			Created:  now,
			Modified: now,
			Tags:     tags,
		},
		Content: content,
	}
}

func (n *Note) SetFilePath(path string) {
	n.FilePath = path
}

func (n *Note) UpdateContent(content string) {
	n.Content = content
	n.Metadata.Modified = time.Now()
}

func (n *Note) UpdateTags(tags []string) {
	n.Metadata.Tags = tags
	n.Metadata.Modified = time.Now()
}

func (n *Note) ToFileContent() (string, error) {
	n.Metadata.Modified = time.Now()

	yamlData, err := yaml.Marshal(&n.Metadata)
	if err != nil {
		return "", fmt.Errorf("error marshaling metadata: %w", err)
	}

	return fmt.Sprintf("---\n%s---\n\n%s", string(yamlData), n.Content), nil
}

func (n *Note) Save() error {
	content, err := n.ToFileContent()
	if err != nil {
		return err
	}

	return os.WriteFile(n.FilePath, []byte(content), 0644)
}