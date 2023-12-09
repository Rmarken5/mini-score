package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestGameDAOImpl_InsertGame(t *testing.T) {
	var (
		inputGame = Game{
			ID:        "1",
			GameTime:  time.Now(),
			Quarter:   "1",
			GameClock: "15:00",
			AwayTeam:  "PIT",
			HomeTeam:  "CIN",
		}
	)
	testCases := map[string]struct {
		input  Game
		mockDB func(db *sql.DB, sqlMock sqlmock.Sqlmock)
		err    error
	}{
		"should return nil error when successful": {
			input: inputGame,
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(regexp.QuoteMeta("insert into game (id, game_time, quarter, game_clock, away_team, home_team) values ($1, $2, $3, $4, $5, $6);")).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"should return error when unsuccessful": {
			input: inputGame,
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(regexp.QuoteMeta("insert into game (id, game_time, quarter, game_clock, away_team, home_team) values ($1, $2, $3, $4, $5, $6);")).WillReturnError(sql.ErrConnDone)
			},
			err: ErrInsertGame,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Mock the database query
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tc.mockDB(db, mock)
			// Create a new TeamDAOImpl with the mocks database and logger
			dao := &GameDAOImpl{
				logger: zerolog.Nop(),
				db:     sqlx.NewDb(db, "postgres"),
			}
			// Call the function
			err = dao.InsertGame(tc.input)

			// Assert the results
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestGameDAOImpl_GetGames(t *testing.T) {
	endTime := time.Now()
	testCases := map[string]struct {
		start         time.Time
		end           *time.Time
		mockDB        func(db *sql.DB, sqlMock sqlmock.Sqlmock)
		expectedCount int
		expectedErr   error
	}{
		"should get games between times": {
			start: time.Now(),
			end:   &endTime,
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				rows := sqlMock.NewRows([]string{"id", "game_time", "quarter", "game_clock", "away_team", "home_team", "created_at", "updated_at", "deleted_at"})
				rows.AddRow(1, time.Now(), "1", "15:00", uuid.New(), uuid.New(), time.Now(), time.Now(), nil)
				rows.AddRow(1, time.Now(), "1", "15:00", uuid.New(), uuid.New(), time.Now(), time.Now(), nil)
				sqlMock.ExpectQuery(regexp.QuoteMeta("select id, game_time, quarter, game_clock, away_team, home_team, created_at, updated_at, deleted_at from game where deleted_at is null and game_time >= $1 and game_time < $2")).WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedErr:   nil,
		},
		"should get games after start": {
			start: time.Now(),
			end:   nil,
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				rows := sqlMock.NewRows([]string{"id", "game_time", "quarter", "game_clock", "away_team", "home_team", "created_at", "updated_at", "deleted_at"})
				rows.AddRow(1, time.Now(), "1", "15:00", uuid.New(), uuid.New(), time.Now(), time.Now(), nil)
				rows.AddRow(1, time.Now(), "1", "15:00", uuid.New(), uuid.New(), time.Now(), time.Now(), nil)
				sqlMock.ExpectQuery(regexp.QuoteMeta(getGamesStmt)).WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedErr:   nil,
		},
		"should return err no games when no games": {
			start: time.Now(),
			end:   &endTime,
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(regexp.QuoteMeta("select id, game_time, quarter, game_clock, away_team, home_team, created_at, updated_at, deleted_at from game where deleted_at is null and game_time >= $1 and game_time < $2")).WillReturnError(sql.ErrNoRows)
			},
			expectedCount: 0,
			expectedErr:   ErrNoGames,
		},
		"generic error": {
			start: time.Now(),
			end:   &endTime,
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(regexp.QuoteMeta("select id, game_time, quarter, game_clock, away_team, home_team, created_at, updated_at, deleted_at from game where deleted_at is null and game_time >= $1 and game_time < $2")).WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectedErr:   ErrSqlError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Mock the database query
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tc.mockDB(db, mock)
			// Create a new TeamDAOImpl with the mocks database and logger
			dao := &GameDAOImpl{
				logger: zerolog.Nop(),
				db:     sqlx.NewDb(db, "postgres"),
			}
			// Call the function
			games, err := dao.GetGames(tc.start, tc.end)

			// Assert the results
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Len(t, games, tc.expectedCount)
		})
	}
}

func TestGameDAOImpl_GetGame(t *testing.T) {
	var (
		haveHomeTeam = uuid.New()
		haveAwayTeam = uuid.New()
		haveId       = "1"
		haveGameTime = time.Now()
	)
	testCases := map[string]struct {
		mockDB       func(db *sql.DB, sqlMock sqlmock.Sqlmock)
		expectedGame Game
		expectedErr  error
	}{
		"should get game for id": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				rows := sqlMock.NewRows([]string{"id", "game_time", "quarter", "game_clock", "away_team", "home_team", "created_at", "updated_at", "deleted_at"})
				rows.AddRow(haveId, haveGameTime, "1", "15:00", haveAwayTeam, haveHomeTeam, haveGameTime, haveGameTime, nil)
				sqlMock.ExpectQuery(regexp.QuoteMeta(getGameStmt)).WillReturnRows(rows)
			},
			expectedGame: Game{
				ID:        haveId,
				GameTime:  haveGameTime,
				Quarter:   "1",
				GameClock: "15:00",
				AwayTeam:  haveAwayTeam.String(),
				HomeTeam:  haveHomeTeam.String(),
				CreatedAt: haveGameTime,
				UpdatedAt: haveGameTime,
				DeletedAt: nil,
			},
			expectedErr: nil,
		},
		"should return err no game when no game": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(getGameStmt)).WillReturnError(sql.ErrNoRows)
			},
			expectedGame: Game{},
			expectedErr:  ErrNoGame,
		},
		"generic error": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(getGameStmt)).WillReturnError(sql.ErrConnDone)
			},
			expectedGame: Game{},
			expectedErr:  ErrSqlError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Mock the database query
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tc.mockDB(db, mock)
			// Create a new TeamDAOImpl with the mocks database and logger
			dao := &GameDAOImpl{
				logger: zerolog.Nop(),
				db:     sqlx.NewDb(db, "postgres"),
			}
			// Call the function
			game, err := dao.GetGame("1")

			// Assert the results
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedGame, game)
		})
	}
}

func TestGameDAOImpl_UpdateQuarterGameClock(t *testing.T) {

	testCases := map[string]struct {
		mockDB func(db *sql.DB, sqlMock sqlmock.Sqlmock)
		err    error
	}{
		"should return nil error when successful": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(regexp.QuoteMeta(updateQuarterGameClock)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"should return error when unsuccessful": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(regexp.QuoteMeta(updateQuarterGameClock)).WillReturnError(sql.ErrConnDone)
			},
			err: ErrUpdateGameClock,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Mock the database query
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tc.mockDB(db, mock)
			// Create a new TeamDAOImpl with the mocks database and logger
			dao := &GameDAOImpl{
				logger: zerolog.Nop(),
				db:     sqlx.NewDb(db, "postgres"),
			}
			// Call the function
			err = dao.UpdateQuarterGameClock("1", "1", uuid.NewString())

			// Assert the results
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestGameDAOImpl_GetGameTeamQuarterScore(t *testing.T) {
	var (
		haveId = "1"
	)
	testCases := map[string]struct {
		mockDB        func(db *sql.DB, sqlMock sqlmock.Sqlmock)
		expectedGames []GameTeamQuarterScore
		expectedErr   error
	}{
		"should get games for dates": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				rows := sqlMock.NewRows([]string{"id", "abbreviation", "quarter", "score"})
				rows.AddRow(haveId, "PIT", "1", "0")
				rows.AddRow(haveId, "PIT", "2", "7")
				rows.AddRow(haveId, "PIT", "3", "3")
				rows.AddRow(haveId, "PIT", "4", "10")
				rows.AddRow(haveId, "HOU", "1", "7")
				rows.AddRow(haveId, "HOU", "2", "3")
				rows.AddRow(haveId, "HOU", "3", "7")
				rows.AddRow(haveId, "HOU", "4", "0")
				stmt := getGameTeamQuarterScoreStmt + " and g.game_time < $2" + getGameTeamQuarterScoreOrderBy
				sqlMock.ExpectQuery(regexp.QuoteMeta(stmt)).WillReturnRows(rows)
			},
			expectedGames: []GameTeamQuarterScore{
				{
					GameID:           haveId,
					TeamAbbreviation: "PIT",
					Quarter:          "1",
					Score:            "0",
				},
				{
					GameID:           haveId,
					TeamAbbreviation: "PIT",
					Quarter:          "2",
					Score:            "7",
				},
				{
					GameID:           haveId,
					TeamAbbreviation: "PIT",
					Quarter:          "3",
					Score:            "3",
				},
				{
					GameID:           haveId,
					TeamAbbreviation: "PIT",
					Quarter:          "4",
					Score:            "10",
				},
				{
					GameID:           haveId,
					TeamAbbreviation: "HOU",
					Quarter:          "1",
					Score:            "7",
				},
				{
					GameID:           haveId,
					TeamAbbreviation: "HOU",
					Quarter:          "2",
					Score:            "3",
				},
				{
					GameID:           haveId,
					TeamAbbreviation: "HOU",
					Quarter:          "3",
					Score:            "7",
				},
				{
					GameID:           haveId,
					TeamAbbreviation: "HOU",
					Quarter:          "4",
					Score:            "0",
				},
			},
			expectedErr: nil,
		},
		"should return err no games when no game": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				stmt := getGameTeamQuarterScoreStmt + " and g.game_time < $2" + getGameTeamQuarterScoreOrderBy
				sqlMock.ExpectQuery(regexp.QuoteMeta(stmt)).WillReturnError(sql.ErrNoRows)
			},
			expectedGames: []GameTeamQuarterScore(nil),
			expectedErr:   ErrNoGames,
		},
		"generic error": {
			mockDB: func(db *sql.DB, sqlMock sqlmock.Sqlmock) {
				stmt := getGameTeamQuarterScoreStmt + " and g.game_time < $2" + getGameTeamQuarterScoreOrderBy
				sqlMock.ExpectQuery(regexp.QuoteMeta(stmt)).WillReturnError(sql.ErrConnDone)
			},
			expectedGames: []GameTeamQuarterScore(nil),
			expectedErr:   ErrSqlError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Mock the database query
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tc.mockDB(db, mock)
			// Create a new TeamDAOImpl with the mocks database and logger
			dao := &GameDAOImpl{
				logger: zerolog.Nop(),
				db:     sqlx.NewDb(db, "postgres"),
			}
			end := time.Now()
			// Call the function
			games, err := dao.GetGameTeamQuarterScore(time.Now(), &end)

			// Assert the results
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedGames, games)
		})
	}

}
