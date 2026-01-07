package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/shreeram-hegde/go-url-shortener/internal/model"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func migrate(db *sql.DB) error { // initiates the server if it doesn't exist
	// schema is just structed string
	schema := `
	CREATE TABLE IF NOT EXISTS urls (
	code TEXT PRIMARY KEY,
	long_url TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	expires_at TIMESTAMP NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_expires_at ON urls(expires_at);
	`
	//db.Exec() just runs the query in the database db

	_, err := db.Exec(schema)
	return err
}

func (s *SQLiteStore) Save(u model.URL) error {
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO urls(code, long_url, created_at, expires_at)
		VALUES(?, ?, ?, ?)`,
		u.Code, u.LongURL, u.CreatedAt, u.ExpiresAt,
	)
	return err
}

func (s *SQLiteStore) Get(code string) (model.URL, error) {
	row := s.db.QueryRow(
		`SELECT code, long_url, created_at, expires_at FROM urls WHERE code = ?`,
		code,
	)

	var u model.URL

	err := row.Scan(&u.Code, &u.LongURL, &u.CreatedAt, &u.ExpiresAt)
	if err == sql.ErrNoRows {
		fmt.Println("not found in db get")
		return model.URL{}, ErrNotFound
	}
	if err != nil {
		fmt.Println(err, "this is in db get")
		return model.URL{}, err
	}

	return u, nil
}

func (s *SQLiteStore) Delete(code string) error {
	_, err := s.db.Exec(`DELETE FROM urls WHERE code = ?`, code)
	return err
}

func (s *SQLiteStore) DeleteExpired(now time.Time) error {
	_, err := s.db.Exec(`DELETE FROM urls WHERE expires_at < ?`, now)
	return err
}
