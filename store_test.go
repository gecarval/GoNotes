package main

import (
	"strings"
	"testing"
)

// newTestStore opens a Store in a fresh temp directory, because Init
// always uses ./notes.db relative to the working directory.
func newTestStore(t *testing.T) *Store {
	t.Helper()
	t.Chdir(t.TempDir())
	s := &Store{}
	if err := s.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	t.Cleanup(func() { s.conn.Close() })
	return s
}

func TestGetNotes(t *testing.T) {
	tests := []struct {
		name  string
		notes []Note
	}{
		{"empty table", nil},
		{"one note", []Note{{1, "title", "body"}}},
		{"empty strings and multi-byte text", []Note{
			{1, "", ""},
			{2, "日本語 ✓", "line one\nline two 🐐"},
		}},
		{"large body", []Note{{1, "big", strings.Repeat("x", 1<<20)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStore(t)
			for _, n := range tt.notes {
				if _, err := s.conn.Exec(`INSERT INTO notes (id, title, body) VALUES (?, ?, ?)`,
					n.id, n.title, n.body); err != nil {
					t.Fatalf("insert: %v", err)
				}
			}
			got, err := s.GetNotes()
			if err != nil {
				t.Fatalf("GetNotes: %v", err)
			}
			if len(got) != len(tt.notes) {
				t.Fatalf("got %d notes, want %d", len(got), len(tt.notes))
			}
			for i := range got {
				if got[i] != tt.notes[i] {
					t.Errorf("note %d = %+v, want %+v", i, got[i], tt.notes[i])
				}
			}
		})
	}
}

// TestGetNotesQueryError covers the error path: a failing query must
// return an error instead of panicking on a nil *sql.Rows.
func TestGetNotesQueryError(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.conn.Exec(`DROP TABLE notes`); err != nil {
		t.Fatalf("drop: %v", err)
	}
	if _, err := s.GetNotes(); err == nil {
		t.Fatal("expected an error for a missing table")
	}
}

// TestGetNotesScanError covers a row that cannot be scanned into a Note:
// the error must surface rather than yield a zero-valued note.
func TestGetNotesScanError(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.conn.Exec(`DROP TABLE notes`); err != nil {
		t.Fatalf("drop: %v", err)
	}
	if _, err := s.conn.Exec(`CREATE TABLE notes (id text, title text, body text)`); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := s.conn.Exec(`INSERT INTO notes VALUES ('not-a-number', 't', 'b')`); err != nil {
		t.Fatalf("insert: %v", err)
	}
	if _, err := s.GetNotes(); err == nil {
		t.Fatal("expected a scan error")
	}
}
