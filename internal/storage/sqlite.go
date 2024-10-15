package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db *sql.DB
}

type RepositoryRecord struct {
	URL         string
	LastUpdated time.Time
}

func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := createTable(db); err != nil {
		return nil, err
	}

	return &SQLiteDB{db: db}, nil
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS repositories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL UNIQUE,
		last_updated DATETIME NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func (s *SQLiteDB) UpdateRepository(url string) error {
	query := `
	INSERT INTO repositories (url, last_updated)
	VALUES (?, ?)
	ON CONFLICT(url) DO UPDATE SET last_updated=excluded.last_updated;`

	_, err := s.db.Exec(query, url, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update repository: %w", err)
	}
	return nil
}

func (s *SQLiteDB) GetOldRepositories() ([]RepositoryRecord, error) {
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	query := `
	SELECT url, last_updated FROM repositories
	WHERE last_updated < ?
	ORDER BY last_updated ASC;`

	rows, err := s.db.Query(query, oneWeekAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to query repositories: %w", err)
	}
	defer rows.Close()

	var records []RepositoryRecord
	for rows.Next() {
		var record RepositoryRecord
		if err := rows.Scan(&record.URL, &record.LastUpdated); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}
