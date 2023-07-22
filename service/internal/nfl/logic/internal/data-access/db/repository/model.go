package repository

import (
	"github.com/google/uuid"
	"time"
)

type Team struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Abbreviation string     `json:"abbreviation" db:"abbreviation"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at,omitempty"`
}

type Game struct {
	ID        string     `json:"id" db:"id"`
	GameTime  time.Time  `json:"game_time" db:"game_time"`
	Quarter   string     `json:"quarter" db:"quarter"`
	GameClock string     `json:"game_clock" db:"game_clock"`
	AwayTeam  string     `json:"away_team" db:"away_team"`
	HomeTeam  string     `json:"home_team" db:"home_team"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at,omitempty"`
}

type GameQuarterScore struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	GameID    string     `json:"game_id" db:"game_id"`
	TeamID    string     `json:"team_id" db:"team_id"`
	Quarter   string     `json:"quarter" db:"quarter"`
	Score     int        `json:"score" db:"score"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at,omitempty"`
}

type GameScoreDTO struct {
	GameID     string `json:"game_id" db:"game_id"`
	TeamScores []ScoreDTO
}

type ScoreDTO struct {
	TeamAbbreviation string
	Score            []string // Each index represents a quarter
}

type GameTeamQuarterScore struct {
	GameID           string `db:"id"`
	TeamAbbreviation string `db:"abbreviation"`
	Quarter          string `db:"quarter"`
	Score            string `db:"score"`
}
