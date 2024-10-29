package create_table

import (
	"database/sql"
	"fmt"
)

// CreateArtistTable creates the Artist table in the database.
func CreateArtistTable(db *sql.DB) error {
	createQuery := `CREATE TABLE IF NOT EXISTS Artist (
		id SERIAL PRIMARY KEY,
		name varchar(50) NOT NULL,
		age INTEGER,
		sex varchar(10),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("error while creating the table: %w", err)
	}

	fmt.Println("Successfully created the Artist table")
	return nil
}
