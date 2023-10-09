package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type TeamDAOImpl struct {
	logger zerolog.Logger
	db     *sqlx.DB
}

func NewTeamDAOImpl(logger zerolog.Logger, db *sqlx.DB) *TeamDAOImpl {
	return &TeamDAOImpl{
		logger: logger.With().Str("repo", "team").Logger(),
		db:     db,
	}
}

const getTeamByAbvStatement = `SELECT ID, NAME, ABBREVIATION, CREATED_AT, UPDATED_AT, DELETED_AT FROM TEAM WHERE ABBREVIATION = $1 AND DELETED_AT IS NULL`

func (t *TeamDAOImpl) GetTeamByAbv(abbv string) (*Team, error) {
	t.logger.Printf("getting team for: %s", abbv)
	team := &Team{}
	err := t.db.Get(team, getTeamByAbvStatement, abbv)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			t.logger.Printf("No team for %s, %s", abbv, err)
			return nil, ErrNoTeam
		default:
			t.logger.Printf("sql error: %s", err)
			return nil, ErrSqlError
		}
	}

	return team, nil
}

const getAllTeamsStmt = `SELECT ID, NAME, ABBREVIATION, CREATED_AT, UPDATED_AT, DELETED_AT FROM TEAM WHERE DELETED_AT IS NULL `

func (t *TeamDAOImpl) GetAllTeams() ([]*Team, error) {
	t.logger.Printf("getting all teams")
	var teams []*Team
	err := t.db.Select(&teams, getAllTeamsStmt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			t.logger.Printf("No teams %s", err)
			return nil, ErrNoTeam
		default:
			t.logger.Printf("sql error: %s", err)
			return nil, ErrSqlError
		}
	}

	return teams, nil
}
