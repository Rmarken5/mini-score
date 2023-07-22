package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rmarken5/mini-score/service/cmd/internal"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/scheduler"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/scheduler/controller"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	logger := createLogger()

	db := internal.MustConnectDatabase(logger)
	defer db.Close()

	sch := createScheduler(logger, db)

	sch.Run()

}

func createLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.With().Str("service", "scheduleRunner").Logger()
	return logger
}

func createScheduler(logger zerolog.Logger, db *sqlx.DB) *scheduler.Scheduler {

	ctrl := controller.NewLogic(logger, db)
	sch := scheduler.New(logger, ctrl)
	return sch
}
