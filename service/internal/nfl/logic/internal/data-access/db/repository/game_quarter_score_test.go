package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestGameQuarterScoreDAOImpl_InsertQuarterScore(t *testing.T) {

	testCases := map[string]struct {
		mockDB func(db *sql.DB, sqlMock sqlmock.Sqlmock)
		err    error
	}{
		"should return nil error when successful": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(regexp.QuoteMeta("insert into game_quarter_score (game_id, team_id, quarter, score, created_at, updated_at, deleted_at) values ($1, (select id from team where abbreviation = $2), $3, $4, $5, $6, $7);")).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"should return error when unsuccessful": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(regexp.QuoteMeta("insert into game_quarter_score (game_id, team_id, quarter, score, created_at, updated_at, deleted_at) values ($1, (select id from team where abbreviation = $2), $3, $4, $5, $6, $7);")).WillReturnError(sql.ErrConnDone)
			},
			err: ErrInsertQuarterScore,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Mock the database query
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tc.mockDB(db, mock)
			// Create a new TeamDAOImpl with the mocks database and logger
			dbx := sqlx.NewDb(db, "postgres")
			dao := &GameQuarterScoreDAOImpl{
				logger: zerolog.Nop(),
				db:     dbx,
			}

			// Call the function
			err = dao.InsertQuarterScore(GameQuarterScore{
				GameID:    "0",
				TeamID:    uuid.NewString(),
				Quarter:   "1",
				Score:     0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: nil,
			})

			// Assert the results
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestGameQuarterScoreDAOImpl_UpdateQuarterScore(t *testing.T) {

	testCases := map[string]struct {
		mockDB func(db *sql.DB, sqlMock sqlmock.Sqlmock)
		err    error
	}{
		"should return nil error when successful": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(regexp.QuoteMeta(updateScoreStmt)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"should return error when unsuccessful": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(updateScoreStmt).WillReturnError(sql.ErrConnDone)
			},
			err: ErrUpdateQuarterScore,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Mock the database query
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tc.mockDB(db, mock)
			// Create a new TeamDAOImpl with the mocks database and logger
			dao := &GameQuarterScoreDAOImpl{
				logger: zerolog.Nop(),
				db:     sqlx.NewDb(db, "postgres"),
			}
			// Call the function
			err = dao.UpdateQuarterScore(21, uuid.NewString(), uuid.NewString(), "1")

			// Assert the results
			assert.ErrorIs(t, err, tc.err)
		})
	}
}
