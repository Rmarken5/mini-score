package internal

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"os"
)

func MustConnectDatabase(logger zerolog.Logger) *sqlx.DB {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DATABASE")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	sslMode := os.Getenv("POSTGRES_SSL_MODE")
	options := os.Getenv("POSTGRES_OPTION")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s options=%s", user, password, database, host, port, sslMode, options)
	// Create the connection pool
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("error opening database connection")
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		logger.Fatal().Err(err).Msg("Error testing database connection")
		panic(err)
	}

	return sqlx.NewDb(db, "postgres")
}
