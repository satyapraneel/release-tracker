package schedular

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/release-trackers/gin/cmd"
)

type Application struct {
	*cmd.Application
}

func Run(app *cmd.Application) {
	application := &Application{app}
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().Do(application.ReleaseDateReminder)
	s.StartAsync()
}
