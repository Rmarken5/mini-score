package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTeamDAOImpl_GetTeamByAbv(t *testing.T) {
	var (
		teamID    = uuid.New()
		createdAt = time.Now()
	)
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	createdAt.In(loc)

	// Define test cases using a map
	testCases := map[string]struct {
		input    string
		mockDB   func(db *sql.DB, sqlMock sqlmock.Sqlmock)
		expected *Team
		err      error
	}{
		"success": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				rows := sqlMock.NewRows([]string{"id", "name", "abbreviation", "created_at", "updated_at"})
				rows.AddRow(teamID, "Test Team", "ABC", createdAt, createdAt)
				sqlMock.ExpectQuery(regexp.QuoteMeta(getTeamByAbvStatement)).WillReturnRows(rows)
			},
			input: "ABC",
			expected: &Team{
				ID:           teamID,
				Name:         "Test Team",
				Abbreviation: "ABC",
				CreatedAt:    createdAt,
				UpdatedAt:    createdAt,
				DeletedAt:    nil,
			},
			err: nil,
		},
		"no team found": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(getTeamByAbvStatement)).WillReturnError(sql.ErrNoRows)
			},
			input:    "DEF",
			expected: nil,
			err:      ErrNoTeam,
		},
		"sql error": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(getTeamByAbvStatement)).WillReturnError(errors.New("whoops"))
			},
			input:    "GHI",
			expected: nil,
			err:      ErrSqlError,
		},
	}

	// Iterate over the test cases and run each one
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			// Mock the database query
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tc.mockDB(db, mock)
			// Create a new TeamDAOImpl with the mocks database and logger
			dao := &TeamDAOImpl{
				logger: zerolog.Nop(),
				db:     sqlx.NewDb(db, "postgres"),
			}
			// Call the function
			team, err := dao.GetTeamByAbv(tc.input)

			// Assert the results
			assert.Equal(t, tc.expected, team)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestTeamDAOImpl_GetAllTeams(t *testing.T) {

	// Define test cases using a map
	testCases := map[string]struct {
		mockDB        func(db *sql.DB, sqlMock sqlmock.Sqlmock)
		expectedCount int
		err           error
	}{
		"should get teams from database": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				rows := sqlMock.NewRows([]string{"id", "name", "abbreviation", "created_at", "updated_at"})
				rows.AddRow(uuid.New(), "Test Team", "ABC", time.Now(), time.Now())
				rows.AddRow(uuid.New(), "Test Team", "ABC", time.Now(), time.Now())
				sqlMock.ExpectQuery(regexp.QuoteMeta(getAllTeamsStmt)).WillReturnRows(rows)
			},
			expectedCount: 2,
			err:           nil,
		},
		"no team found": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(getAllTeamsStmt)).WillReturnError(sql.ErrNoRows)
			},
			expectedCount: 0,
			err:           ErrNoTeam,
		},
		"sql error": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(getAllTeamsStmt)).WillReturnError(errors.New("whoops"))
			},
			expectedCount: 0,
			err:           ErrSqlError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			// Mock the database query
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tc.mockDB(db, mock)
			// Create a new TeamDAOImpl with the mocks database and logger
			dao := &TeamDAOImpl{
				logger: zerolog.Nop(),
				db:     sqlx.NewDb(db, "postgres"),
			}
			// Call the function
			teams, err := dao.GetAllTeams()

			// Assert the results
			assert.Len(t, teams, tc.expectedCount)
			assert.Equal(t, tc.err, err)
		})
	}

}
