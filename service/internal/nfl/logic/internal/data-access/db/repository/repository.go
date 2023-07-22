package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"time"
)

type (
	TeamDAO interface {
		GetTeamByAbv(abbv string) (*Team, error)
		GetAllTeams() ([]*Team, error)
	}

	GameDAO interface {
		InsertGame(game Game) error
		GetGames(start time.Time, end *time.Time) ([]Game, error)
		GetGamesWithTeamAbv(start time.Time, end *time.Time) ([]Game, error)
		GetGame(gameID string) (Game, error)
		UpdateQuarterGameClock(gameID string, quarter string, gameClock string) error
		UpdateGameTime(gameID string, gameClock time.Time) error
		UpdateGameClock(gameID string, gameClock string) error

		GetGameTeamQuarterScore(start time.Time, end *time.Time) ([]GameTeamQuarterScore, error)
	}

	GameQuarterScoreDAO interface {
		GetQuarterScoreBy(gameID string, teamAbv string, quarter string) (GameQuarterScore, error)
		InsertQuarterScore(quarterScore GameQuarterScore) error
		UpdateQuarterScore(score int, gameID string, teamAbv string, quarter string) error
	}

	Repository interface {
		TeamDAO
		GameDAO
		GameQuarterScoreDAO
	}

	RepositoryImpl struct {
		TeamDAO
		GameDAO
		GameQuarterScoreDAO
	}
)

func NewRepository(logger zerolog.Logger, db *sqlx.DB) *RepositoryImpl {
	teamDAO := NewTeamDAOImpl(logger, db)
	gameDAO := NewGameDAOImpl(logger, db)
	gameQuarterScoreDAO := NewGameQuarterScoreDAOImpl(logger, db)
	return &RepositoryImpl{
		TeamDAO:             teamDAO,
		GameDAO:             gameDAO,
		GameQuarterScoreDAO: gameQuarterScoreDAO,
	}
}
