package analytics

import (
	"database/sql"
)

// GetTotalEvents asks the Postgres vault for the row count
func GetTotalEvents(db *sql.DB) (int, error) {
	var count int
	
	// Senior Tip: We use QueryRow when we expect exactly one result
	err := db.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
	
	if err != nil {
		return 0, err
	}
	
	return count, nil
}