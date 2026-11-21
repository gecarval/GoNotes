package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Note is a single note as stored in the notes table.
type Note struct {
	id    int64
	title string
	body  string
}

// Store persists notes in a local SQLite database.
type Store struct {
	conn *sql.DB
}

// Init opens ./notes.db and creates the notes table if it is missing.
func (s *Store) Init() (err error) {
	if s.conn, err = sql.Open("sqlite3", "./notes.db"); err != nil {
		return err
	}
	createTableStatement := `CREATE TABLE IF NOT EXISTS notes (
		id integer not null primary key,
		title text not null,
		body text not null
	);`
	if _, err := s.conn.Exec(createTableStatement); err != nil {
		return err
	}
	return nil
}

// GetNotes returns every note in the store.
func (s *Store) GetNotes() (notes []Note, err error) {
	getNotesStatement := `SELECT * FROM notes`
	rows, err := s.conn.Query(getNotesStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var note Note
		if err = rows.Scan(&note.id, &note.title, &note.body); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	// rows.Err reports iteration errors only once Next has returned false,
	// so it must be checked after the loop, not inside it.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}
