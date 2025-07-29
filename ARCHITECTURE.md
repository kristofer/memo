# Memo Application Architecture

This document outlines the architectural structure of the Memo CLI application (a version of Future-Proof-Notes) after refactoring to implement clean separation of concerns and the Command Pattern.

## Overall Architecture

```mermaid
graph TB
    %% Entry Point
    Main[main.go<br/>Entry Point] --> App[cmd/App<br/>Application Controller]
    
    %% Command Layer
    App --> CR[cmd/CommandRegistry<br/>Command Management]
    CR --> CC[CreateCommand]
    CR --> LC[ListCommand]  
    CR --> RC[ReadCommand]
    CR --> EC[EditCommand]
    CR --> DC[DeleteCommand]
    CR --> SC[SearchCommand]
    CR --> STC[StatsCommand]
    CR --> HC[HelpCommand]
    
    %% Command Context
    CC --> CTX[CommandContext<br/>Shared State]
    LC --> CTX
    RC --> CTX
    EC --> CTX
    DC --> CTX
    SC --> CTX
    STC --> CTX
    HC --> CTX
    
    %% Internal Layers
    CTX --> Storage[internal/storage<br/>FileStorage]
    CC --> UI[internal/ui<br/>User Interface]
    LC --> UI
    RC --> UI
    EC --> UI
    DC --> UI
    SC --> UI
    STC --> UI
    HC --> UI
    
    %% Domain Layer
    Storage --> NotePackage[internal/note<br/>Domain Models]
    UI --> NotePackage
    
    %% External Dependencies
    Storage --> FS[File System<br/>.memo-notes/]
    UI --> Console[Console I/O<br/>stdin/stdout]
    Note --> YAML[YAML Library<br/>gopkg.in/yaml.v3]
    
    %% Styling
    classDef entryPoint fill:#e1f5fe
    classDef cmdLayer fill:#f3e5f5
    classDef internal fill:#e8f5e8
    classDef domain fill:#fff3e0
    classDef external fill:#ffebee
    
    class Main entryPoint
    class App,CR,CC,LC,RC,EC,DC,SC,STC,HC,CTX cmdLayer
    class Storage,UI internal
    class NotePackage domain
    class FS,Console,YAML external
```

## Package Structure

```mermaid
graph LR
    subgraph "Application Entry"
        M[main.go]
    end
    
    subgraph "Command Layer (cmd/)"
        A[App]
        CI[Command Interface]
        CC[Command Implementations]
        CTX[CommandContext]
    end
    
    subgraph "Internal Packages (internal/)"
        subgraph "Domain (note/)"
            NoteModel[Note]
            NM[Metadata]
        end
        
        subgraph "Storage (storage/)"
            FS[FileStorage]
        end
        
        subgraph "UI (ui/)"
            UI[User Interface]
        end
    end
    
    subgraph "External Dependencies"
        YAML[YAML Parser]
        FileSystem[File System]
        Console[Console I/O]
    end
    
    %% Dependencies
    M --> A
    A --> CI
    CI --> CC
    CC --> CTX
    CTX --> FS
    CC --> UI
    FS --> NoteModel
    UI --> NoteModel
    FS --> FileSystem
    UI --> Console
    NoteModel --> YAML
```

## Dependency Flow

```mermaid
graph TD
    %% Layer definitions
    subgraph "Presentation Layer"
        CLI[CLI Commands]
    end
    
    subgraph "Application Layer"
        APP[App Controller]
        CMD[Command Pattern]
    end
    
    subgraph "Domain Layer"
        NOTE[Note Domain]
    end
    
    subgraph "Infrastructure Layer"
        STORAGE[File Storage]
        UI[User Interface]
    end
    
    subgraph "External"
        FILES[File System]
        CONSOLE[Console]
    end
    
    %% Dependencies (top to bottom)
    CLI --> APP
    APP --> CMD
    CMD --> NOTE
    CMD --> STORAGE
    CMD --> UI
    STORAGE --> FILES
    UI --> CONSOLE
    
    %% Clean Architecture Rules
    NOTE -.->|"No dependencies on outer layers"| NOTE
```

## Command Pattern Implementation

```mermaid
classDiagram
    class Command {
        <<interface>>
        +Execute(args []string) error
    }
    
    class CommandContext {
        +Storage *FileStorage
        +CurrentListing []*Note
        +SetCurrentListing(notes []*Note)
        +GetCurrentListing() []*Note
    }
    
    class App {
        -ctx *CommandContext
        -commands map[string]Command
        +NewApp() *App
        +registerCommands()
        +Run()
    }
    
    class CreateCommand {
        -ctx *CommandContext
        +Execute(args []string) error
    }
    
    class ListCommand {
        -ctx *CommandContext  
        +Execute(args []string) error
    }
    
    class ReadCommand {
        -ctx *CommandContext
        +Execute(args []string) error
        +resolveNoteID(identifier string) (string, error)
    }
    
    %% Relationships
    Command <|.. CreateCommand
    Command <|.. ListCommand
    Command <|.. ReadCommand
    App --> CommandContext
    App --> Command
    CreateCommand --> CommandContext
    ListCommand --> CommandContext
    ReadCommand --> CommandContext
```

## Domain Model

```mermaid
classDiagram
    class Note {
        +Metadata Metadata
        +Content string
        +FilePath string
        +New(title, content string, tags []string) *Note
        +SetFilePath(path string)
        +UpdateContent(content string)
        +UpdateTags(tags []string)
        +ToFileContent() (string, error)
        +Save() error
    }
    
    class Metadata {
        +Title string
        +Created time.Time
        +Modified time.Time
        +Tags []string
        +Author string
        +Status string
        +Priority int
    }
    
    class FileStorage {
        -notesDir string
        -noteExtension string
        +NewFileStorage() *FileStorage
        +SaveNote(note *Note) error
        +GetAllNotes() []*Note, error
        +FindNoteByID(noteID string) (*Note, error)
        +DeleteNote(noteID string) error
        +SearchNotes(query string) ([]*Note, error)
        +FilterNotesByTag(tag string) ([]*Note, error)
    }
    
    Note --> Metadata
    FileStorage --> Note
```

## Data Flow

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant App
    participant Command
    participant Storage
    participant NoteModel as Note
    participant FileSystem
    
    User->>CLI: memo create
    CLI->>App: Run()
    App->>App: Parse command line
    App->>Command: CreateCommand.Execute()
    Command->>Command: Prompt for input
    Command->>NoteModel: New(title, content, tags)
    NoteModel-->>Command: *Note
    Command->>Storage: SaveNote(note)
    Storage->>NoteModel: ToFileContent()
    NoteModel-->>Storage: YAML + content
    Storage->>FileSystem: WriteFile()
    FileSystem-->>Storage: success
    Storage-->>Command: nil
    Command-->>App: nil
    App-->>CLI: success
    CLI-->>User: "Note created successfully"
```

## Key Architectural Benefits

### 1. **Separation of Concerns**
- **Domain Logic**: Pure business rules in `internal/note`
- **Storage Logic**: File operations in `internal/storage`  
- **UI Logic**: User interaction in `internal/ui`
- **Command Logic**: CLI handling in `cmd`

### 2. **Command Pattern Benefits**
- **Extensibility**: Easy to add new commands
- **Testability**: Each command can be unit tested
- **Maintainability**: Single responsibility per command
- **Consistency**: All commands follow same interface

### 3. **Clean Architecture Principles**
- **Dependency Inversion**: Outer layers depend on inner layers
- **Interface Segregation**: Small, focused interfaces
- **Single Responsibility**: Each package has one reason to change
- **Open/Closed**: Open for extension, closed for modification

### 4. **Future-Proofing**
- **Plugin Architecture**: Commands can be dynamically registered
- **Storage Abstraction**: Easy to swap file storage for database
- **UI Abstraction**: Can support web UI or GUI in future
- **Domain Purity**: Business logic independent of infrastructure

## Package Responsibilities

| Package | Responsibility | Dependencies |
|---------|---------------|--------------|
| `main` | Application entry point | `cmd` |
| `cmd` | CLI command handling & routing | `internal/*` |
| `internal/note` | Domain models & business logic | Standard library, YAML |
| `internal/storage` | Data persistence operations | `internal/note` |
| `internal/ui` | User interface & interaction | `internal/note` |

This architecture makes the codebase beginner-friendly while maintaining professional standards for scalability and maintainability.