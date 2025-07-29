package cmd

import (
	"fmt"
	"os"

	"memo/internal/storage"
	"memo/internal/ui"
)

type App struct {
	ctx      *CommandContext
	commands map[string]Command
}

func NewApp() *App {
	ctx := &CommandContext{
		Storage: storage.NewFileStorage(),
	}

	app := &App{
		ctx:      ctx,
		commands: make(map[string]Command),
	}

	// Register all commands
	app.registerCommands()
	
	return app
}

func (app *App) registerCommands() {
	app.commands["create"] = NewCreateCommand(app.ctx)
	app.commands["list"] = NewListCommand(app.ctx)
	app.commands["read"] = NewReadCommand(app.ctx)
	app.commands["edit"] = NewEditCommand(app.ctx)
	app.commands["delete"] = NewDeleteCommand(app.ctx)
	app.commands["search"] = NewSearchCommand(app.ctx)
	app.commands["stats"] = NewStatsCommand(app.ctx)
	app.commands["help"] = NewHelpCommand(app.ctx)
	app.commands["--help"] = NewHelpCommand(app.ctx)
	app.commands["-h"] = NewHelpCommand(app.ctx)
}

func (app *App) Run() {
	if len(os.Args) < 2 {
		ui.PrintHelp()
		return
	}

	commandName := os.Args[1]
	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	command, exists := app.commands[commandName]
	if !exists {
		fmt.Printf("Unknown command: %s\n", commandName)
		ui.PrintHelp()
		return
	}

	err := command.Execute(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}