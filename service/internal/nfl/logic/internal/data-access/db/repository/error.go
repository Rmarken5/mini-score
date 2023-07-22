package repository

import "errors"

var (
	ErrNoTeam          = errors.New("no team returned from database")
	ErrSqlError        = errors.New("sql database error")
	ErrInsertGame      = errors.New("error inserting game into database")
	ErrUpdateGameClock = errors.New("error updating game clock in database")
	ErrNoGames         = errors.New("no games returned from database")
	ErrNoGame          = errors.New("no game returned from database")

	ErrNoQuarterScore = errors.New("no quarter score returned from database")

	ErrInsertQuarterScore = errors.New("error inserting quarter score into database")
	ErrUpdateQuarterScore = errors.New("error updating quarter score into database")
)
