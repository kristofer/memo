# Memo - a future-proof Personal Notes Manager

Written in Go, as an example. 

A multi-phase educational project for managing personal notes with structured metadata.

## Overview

Personal Notes Manager is a text-based note-taking system that stores notes as UTF-8 text files with YAML headers for metadata. This project serves as an educational tool for students to learn file manipulation, parsing, and progressively more advanced application architectures.

https://zcw.guru/kristofer/javapygitignore

## Project Phases

This project is designed to be implemented in three distinct phases:

### Phase 1: Command Line Interface (CLI)
- Basic note creation, reading, updating, and deletion via command line
- YAML header parsing and manipulation
- Note searching and filtering capabilities
- Simple formatting options

### Phase 2: Graphical User Interface (GUI)
- Desktop application with intuitive UI
- Rich text editing features
- Visual organization of notes (folders/tags)
- Enhanced search functionality with highlighting
- Export options (PDF, HTML, etc.)

### Phase 3: REST Server with JavaScript Frontend
- Web-based access to notes
- Multi-user support with authentication
- Real-time collaborative editing
- Mobile-responsive design
- API documentation for potential extensions

## File Structure

Notes are stored as plain text files with a `.note` extension (or any extension you prefer). Each note follows this structure:

```
---
title: My Example Note
created: 2025-05-20T10:30:00Z
modified: 2025-05-20T10:45:00Z
tags: [example, documentation]
---

This is the content of my note.
It can contain multiple paragraphs and basic formatting.

You can include lists:
- Item 1
- Item 2
- Item 3

And other simple markup as needed.
```

## YAML Header Specification

The YAML header is delimited by triple dashes (`---`) and contains metadata about the note:

| Field | Description | Type | Required |
|-------|-------------|------|----------|
| title | The title of the note | String | Yes |
| created | Creation timestamp (ISO 8601) | DateTime | Yes |
| modified | Last modification timestamp (ISO 8601) | DateTime | Yes |
| tags | Categories or labels for the note | Array of Strings | No |
| author | Creator of the note | String | No |
| status | Completion or review status | String | No |
| priority | Importance level | Integer (1-5) | No |

Additional custom fields can be added as needed.

## Phase 1: CLI Implementation

### Command Reference

```
notes --help                     # Display help information
notes create                     # Create a new note (opens in default editor)
notes list                       # List all notes
notes list --tag "coursework"    # List notes with specific tag
notes read <note-id>             # Display a specific note
notes edit <note-id>             # Edit a specific note
notes delete <note-id>           # Delete a specific note
notes search "query"             # Search notes for text
notes stats                      # Display statistics about your notes
```

### Technical Requirements

- Java 11+ or Python 3.11+
- YAML parser library
- File system access
- Text editor integration

## Phase 2: GUI Implementation

### Features

- Note browser panel (folder structure or tag-based)
- Rich text editor with markdown support
- Real-time preview
- Drag and drop organization
- Dark/light theme options
- Backup and restore functionality

### Technical Requirements

- Java Swing or TKinter
- HTML/CSS for UI (if using web technologies)
- Local storage management
- Cross-platform compatibility

## Phase 3: REST Server Implementation

### API Endpoints

```
GET    /api/notes                # List all notes
POST   /api/notes                # Create a new note
GET    /api/notes/:id            # Retrieve a specific note
PUT    /api/notes/:id            # Update a specific note
DELETE /api/notes/:id            # Delete a specific note
GET    /api/tags                 # List all tags
GET    /api/notes/tag/:tagid          # List all notes with tag tagid
GET    /api/search?q=query       # Search notes
```

### Technical Requirements

- Express.js/Flask/Django REST framework (python)
- Spring/Springboot
- VanillaJS
- Frontend: React/Vue/Angular
- Authentication system
- Database integration (optional, should still use file system)
- API documentation (Swagger/OpenAPI)

## Learning Objectives

Students will gain experience with:

1. File I/O operations
2. Data parsing and validation
3. Command-line argument processing
4. User interface design principles
5. Client-server architecture
6. REST API development
7. Frontend-backend integration
8. Project planning and progressive enhancement

## Getting Started

### Prerequisites

- [List required software/libraries]
- Basic understanding of Java or Python
- Knowledge of [relevant concepts]

### Installation

```bash
# Clone the repository
# CHOOSE your Project's Name NOW.... _don't use future-proof_.
git clone https://github.com/yourusername/future-proof.git

# Navigate to the project directory or whatever you named it in the clone.
cd `yournotesnamehere`

# Install dependencies
(if you need them)

# Run the application
run.sh
```

## Contributing

This project is designed for educational purposes. Students are encouraged to:

1. Fork the repository
2. Create a feature branch
3. Implement their assigned component
4. Submit a pull request for review

## License

[MIT license]

## Acknowledgments

- Every amazing personal notes project out there.
- Unix file system. Without it, we'd be sunk.
- Markdown: the way techies write.

## Possible Project Names

Possible Names for your version of this project (or make up your own!)

NoteNexus
MindScribe
ThoughtVault
StudyStream
MemoPad Pro
BrainWave Notes
Scholarly Jotter
InsightKeeper
Knowledge Capsule
NoteCraft
MindMapper
StudyPulse
ThoughtHub
NoteWorthy
IdeaTracker
StudySync
RecallPro
ClassCompass
MindfulNotes
LearnLogger

## Finally

_Why does choosing to use text files in a standard directory structure using Markdown as a note format, make the project "future proof"?_

Using text files in a standard directory structure with Markdown as a note format makes a project "future proof" for several compelling reasons:

**Universal Compatibility**
- Plain text files (.txt, .md) are readable by virtually any operating system and text editor, eliminating dependency on proprietary software
- UTF-8 encoding ensures support for multiple languages and special characters across all modern systems

**Longevity and Stability**
- Text files have remained a stable format for decades and will likely be readable for decades to come
- Unlike proprietary formats that can become obsolete when companies discontinue support

**Version Control Friendly**
- Text files integrate seamlessly with version control systems like Git
- Changes can be tracked line by line, enabling precise collaboration and history tracking

**Portable and Accessible**
- Files can be easily transferred between devices and platforms
- Low storage requirements compared to binary formats

**Human Readable**
- Content remains accessible without specialized software
- Even if the application becomes obsolete, notes remain directly readable

**Open Standards**
- Markdown is an open specification with wide industry support
- Not controlled by a single entity that could change or abandon the format

**Flexibility and Extensibility**
- Directory structures can be organized according to user preference
- Metadata in YAML headers can evolve over time without breaking backward compatibility
- New features can be added without rendering old files unusable

**Searchable**
- Plain text is easily searchable using standard system tools or simple scripts
- No need for specialized databases that might become obsolete

**Data Sovereignty**
- Users maintain complete control over their data
- No dependence on cloud services that may change terms or shut down

**Resilient Against Software Evolution**
- Even as the application evolves through its three phases (CLI → GUI → web), the underlying data format remains consistent
- Allows for migration to newer systems without data loss

This approach creates a foundation that can withstand technological change, ensuring that notes remain accessible regardless of future software and hardware developments.
