package app

import (
	"github.com/neutralusername/Systemge/Node"
)

type App struct {
}

func New() Node.Application {
	app := &App{}
	return app
}

func (app *App) GetCommandHandlers() map[string]Node.CommandHandler {
	return map[string]Node.CommandHandler{}
}
