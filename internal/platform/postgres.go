package platform

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // This is the 'driver'
)

// NewPostgresDB opens a connection to the permanent vault
func NewPostgresDB(dsn string) (*sql.DB, error) {
	// 1. Open the connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// 2. Ping to verify the bridge is up
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}