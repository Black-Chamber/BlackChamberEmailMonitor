package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type Storage struct {
	DB *sql.DB
}

func Initialize(path string) *Storage {
	// Check if the database file exists
	dbPath := path + "/BCEM.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("Database file does not exist. Creating: %s", dbPath)
	}

	// Open the database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Ensure the schema is up-to-date
	err = ensureSchema(db)
	if err != nil {
		log.Fatalf("Failed to ensure database schema: %v", err)
	}

	return &Storage{DB: db}
}

// ensureSchema ensures all required tables and schema are present
func ensureSchema(db *sql.DB) error {
	// SQL to create or update tables
	schema := `
	-- Matches Table
	CREATE TABLE IF NOT EXISTS Matches (
		match_id INTEGER PRIMARY KEY AUTOINCREMENT,
		message_id TEXT NOT NULL,
		recipient TEXT NOT NULL,
		service TEXT NOT NULL,
		rule_name TEXT NOT NULL,
		confidence_level INTEGER NOT NULL,
		matched_time DATETIME NOT NULL,
		FOREIGN KEY (message_id) REFERENCES Messages (message_id),
		FOREIGN KEY (rule_name) REFERENCES DetectionRules (rule_name)
	);

	-- Index for faster lookups in Matches
	CREATE INDEX IF NOT EXISTS idx_matches_confidence ON Matches (confidence_level);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema statements: %v", err)
	}

	return nil
}

// InsertMatch inserts a match between a message and a rule
func (s *Storage) InsertMatch(messageID string, recipient string, service string, ruleID string, confidence int) error {
	_, err := s.DB.Exec(`
		INSERT OR IGNORE INTO Matches (message_id, recipient, service, rule_name, confidence_level, matched_time)
		VALUES (?, ?, ?, ?, ?, ?);`,
		messageID, recipient, service, ruleID, confidence, time.Now().UTC(),
	)
	return err
}
