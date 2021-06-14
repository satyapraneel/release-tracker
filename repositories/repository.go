package repositories

import "github.com/release-trackers/gin/cmd"

type App struct {
	*cmd.Application
}

// NewReleaseHandler ..
func NewRepositoryHandler(app *cmd.Application) *App {
	return &App{app}
}
