package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func (s *SQLite) Connect(dsn string) (*sql.DB, error) {
	if s.db != nil {
		return s.db, nil
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	s.db = db
	return s.db, nil
}

func (s *SQLite) Close() {
	if s.db != nil {
		_ = s.db.Close()
	}
}
