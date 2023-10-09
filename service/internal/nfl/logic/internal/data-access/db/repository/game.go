package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"time"
)

type GameDAOImpl struct {
	logger zerolog.Logger
	db     *sqlx.DB
}

func NewGameDAOImpl(logger zerolog.Logger, db *sqlx.DB) *GameDAOImpl {

	return &GameDAOImpl{
		logger: logger.With().Str("repo", "GameDAO").Logger(),
		db:     db,
	}
}

const insertGameStmt = `insert into game (id, game_time, quarter, game_clock, away_team, home_team)
values (:id, :game_time, :quarter, :game_clock, :away_team, :home_team);`

func (g *GameDAOImpl) InsertGame(game Game) error {
	logger := g.logger.With().Str("method", "InsertGame").Logger()
	logger.Info().Msgf("inserting game: %+v", game)

	_, err := g.db.NamedExec(insertGameStmt, &game)
	if err != nil {
		return errors.Join(err, ErrInsertGame)
	}

	return nil
}

const getGamesStmt = "select id, game_time, quarter, game_clock, away_team, home_team, created_at, updated_at, deleted_at from game where deleted_at is null and game_time >= $1"

func (g *GameDAOImpl) GetGames(start time.Time, end *time.Time) ([]Game, error) {
	logger := g.logger.With().Str("method", "GetGames").Logger()
	logger.Info().Msgf("getting games between %s and %s", start, end)

	args := []interface{}{start}
	stmt := getGamesStmt
	if end != nil {
		stmt += " and game_time < $2"
		args = append(args, end)
	}
	games := &[]Game{}
	err := g.db.Select(games, stmt, args...)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logger.Info().Msgf("No game for %s, %s, %s", start, end, err)
			return nil, ErrNoGames
		default:
			logger.Info().Msgf("sql error: %s", err)
			return nil, ErrSqlError
		}
	}

	return *games, nil
}

const GetGamesWithTeamAbvStmt = `SELECT
    g.id,
    t_away.abbreviation AS away_team,
    t_home.abbreviation AS home_team,
    g.game_time, 
    g.game_clock,
    g.quarter
FROM
    game AS g
        INNER JOIN
    team AS t_away ON g.away_team = t_away.id
        INNER JOIN
    team AS t_home ON g.home_team = t_home.id where g.deleted_at is null and g.game_time >= $1`

func (g *GameDAOImpl) GetGamesWithTeamAbv(start time.Time, end *time.Time) ([]Game, error) {
	logger := g.logger.With().Str("method", "GetGamesWithTeamAbv").Logger()
	logger.Info().Msgf("getting games between %s and %s", start, end)

	args := []interface{}{start}
	stmt := GetGamesWithTeamAbvStmt
	if end != nil {
		stmt += " and g.game_time < $2"
		args = append(args, end)
	}
	games := &[]Game{}
	err := g.db.Select(games, stmt, args...)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logger.Info().Msgf("No game for %s, %s, %s", start, end, err)
			return nil, ErrNoGames
		default:
			logger.Info().Msgf("sql error: %s", err)
			return nil, ErrSqlError
		}
	}

	return *games, nil
}

const getGameStmt = "select id, game_time, quarter, game_clock, away_team, home_team, created_at, updated_at, deleted_at from game where id=$1 and deleted_at is null;"

func (g *GameDAOImpl) GetGame(gameID string) (Game, error) {
	logger := g.logger.With().Str("method", "GetGame").Logger()
	logger.Info().Msgf("getting game for gameID %s", gameID)

	var game = &Game{}
	err := g.db.Get(game, getGameStmt, gameID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logger.Info().Msgf("No game for %s, %s", gameID, err)
			return Game{}, ErrNoGame
		default:
			return Game{}, errors.Join(err, ErrSqlError)
		}
	}

	return *game, nil
}

const updateQuarterGameClock = "UPDATE GAME set game_clock = $1, quarter = $2 where id = $3 and deleted_at is null"

func (g *GameDAOImpl) UpdateQuarterGameClock(gameID string, quarter string, gameClock string) error {
	logger := g.logger.With().Str("method", "UpdateQuarterGameClock").Logger()
	logger.Info().Msgf("updating game %s", gameID)

	_, err := g.db.Exec(updateQuarterGameClock, gameClock, quarter, gameID)
	if err != nil {
		logger.Info().Msgf("error updating game %s, %s", gameID, err)
		return errors.Join(err, ErrUpdateGameClock)
	}

	return nil
}

const updateGameClock = "UPDATE GAME set game_clock = $1 where id = $2 and deleted_at is null"

func (g *GameDAOImpl) UpdateGameClock(gameID string, gameClock string) error {
	logger := g.logger.With().Str("method", "UpdateGameClock").Logger()
	logger.Info().Msgf("updating game %s", gameID)

	_, err := g.db.Exec(updateGameClock, gameClock, gameID)
	if err != nil {
		logger.Info().Msgf("error updating game clock %s, %s", gameID, err)
		return errors.Join(err, ErrUpdateGameClock)
	}

	return nil
}

const updateGameTimeStmt = "UPDATE GAME SET game_time=$1 WHERE id=$2"

func (g *GameDAOImpl) UpdateGameTime(gameID string, gameTime time.Time) error {
	logger := g.logger.With().Str("method", "UpdateGameTime").Logger()
	logger.Info().Msgf("updating game %s with gameTime: %v", gameID, gameTime)
	_, err := g.db.Exec(updateGameTimeStmt, gameTime, gameID)
	if err != nil {
		return err
	}
	return nil
}

// language=sql
const getGameTeamQuarterScoreStmt = `select g.id, t.abbreviation, gqs.quarter, gqs.score
from game g
        inner join game_quarter_score gqs on g.id = gqs.game_id
         inner join team t on t.id = gqs.team_id
where g.game_time >= $1
`
const getGameTeamQuarterScoreOrderBy = " order by g.id, g.game_time, gqs.team_id, gqs.quarter"

func (g *GameDAOImpl) GetGameTeamQuarterScore(start time.Time, end *time.Time) ([]GameTeamQuarterScore, error) {

	logger := g.logger.With().Str("method", "GetGameTeamQuarterScore").Logger()
	logger.Info().Msgf("getting games team quarter between %s and %s", start, end)

	args := []interface{}{start}
	stmt := getGameTeamQuarterScoreStmt
	if end != nil {
		stmt += " and g.game_time < $2"
		args = append(args, end)
	}
	stmt += getGameTeamQuarterScoreOrderBy

	var quarterScores []GameTeamQuarterScore
	err := g.db.Select(&quarterScores, stmt, args...)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logger.Info().Msgf("No games for %s, %s, %s", start, end, err)
			return nil, ErrNoGames
		default:
			logger.Info().Msgf("sql error: %s", err)
			return nil, ErrSqlError
		}
	}

	return quarterScores, nil
}
