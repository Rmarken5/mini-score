package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

var _ GameQuarterScoreDAO = &GameQuarterScoreDAOImpl{}

type GameQuarterScoreDAOImpl struct {
	logger zerolog.Logger
	db     *sqlx.DB
}

func NewGameQuarterScoreDAOImpl(logger zerolog.Logger, db *sqlx.DB) *GameQuarterScoreDAOImpl {
	return &GameQuarterScoreDAOImpl{
		logger: logger.With().Str("repo", "GameQuarterScoreDAO").Logger(),
		db:     db,
	}
}

const getQuarterScoreByUniqueStmt = "select id, game_id, team_id, quarter, score, created_at, updated_at, deleted_at from game_quarter_score where game_id = $1 and team_id=(select id from team where abbreviation = $2) and quarter = $3;"

func (g *GameQuarterScoreDAOImpl) GetQuarterScoreBy(gameID string, teamAbv string, quarter string) (GameQuarterScore, error) {
	logger := g.logger.With().Str("method", "InsertQuarterScore").Logger()

	logger.Info().Msgf("getting quarterScore for game: %s, team: %s, quarter: %s", gameID, teamAbv, quarter)
	var gqs = GameQuarterScore{}
	err := g.db.Get(&gqs, getQuarterScoreByUniqueStmt, gameID, teamAbv, quarter)
	if err != nil {
		logger.Error().Err(err).Msg("while getting quarter score")
		switch err {
		case sql.ErrNoRows:
			return GameQuarterScore{}, ErrNoQuarterScore
		default:
			return GameQuarterScore{}, ErrSqlError
		}
	}

	return gqs, err
}

const insertQuarterScoreStmt = "insert into game_quarter_score (game_id, team_id, quarter, score, created_at, updated_at, deleted_at) values (:game_id, (select id from team where abbreviation = :team_id), :quarter, :score, :created_at, :updated_at, :deleted_at);"

func (g *GameQuarterScoreDAOImpl) InsertQuarterScore(quarterScore GameQuarterScore) error {
	logger := g.logger.With().Str("method", "InsertQuarterScore").Logger()
	logger.Info().Msgf("inserting quarterScore: %+v", quarterScore)
	_, err := g.db.NamedExec(insertQuarterScoreStmt, &quarterScore)
	if err != nil {
		return errors.Join(fmt.Errorf("error inserting quarterScore: %+v. %w", quarterScore, err), ErrInsertQuarterScore)
	}

	return nil
}

const updateScoreStmt = "UPDATE GAME_QUARTER_SCORE SET SCORE=$1 WHERE game_id=$2 AND team_id=(select id from team where abbreviation = $3) AND quarter = $4"

func (g *GameQuarterScoreDAOImpl) UpdateQuarterScore(score int, gameID string, teamAbv string, quarter string) error {
	logger := g.logger.With().Str("method", "UpdateQuarterScore").Logger()
	logger.Info().Msgf("updating score GAME: %s, TEAM: %s, QUARTER: %s", gameID, teamAbv, quarter)

	_, err := g.db.Exec(updateScoreStmt, score, gameID, teamAbv, quarter)
	if err != nil {
		return errors.Join(err, ErrUpdateQuarterScore)
	}

	return nil
}
