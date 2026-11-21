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

func TestSaveNote(t *testing.T) {
	tests := []struct {
		name  string
		saves []Note
		want  []Note
	}{
		{"insert with explicit id", []Note{{7, "t", "b"}}, []Note{{7, "t", "b"}}},
		{"update existing id",
			[]Note{{7, "old", "old body"}, {7, "new", "new body"}},
			[]Note{{7, "new", "new body"}}},
		{"two distinct ids",
			[]Note{{1, "a", "x"}, {2, "b", "y"}},
			[]Note{{1, "a", "x"}, {2, "b", "y"}}},
		{"empty and multi-byte strings",
			[]Note{{3, "", ""}, {4, "日本語 ✓", "line\n🐐"}},
			[]Note{{3, "", ""}, {4, "日本語 ✓", "line\n🐐"}}},
		{"large body", []Note{{5, "big", strings.Repeat("x", 1<<20)}},
			[]Note{{5, "big", strings.Repeat("x", 1<<20)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStore(t)
			for _, n := range tt.saves {
				if err := s.SaveNote(n); err != nil {
					t.Fatalf("SaveNote(id %d): %v", n.id, err)
				}
			}
			got, err := s.GetNotes()
			if err != nil {
				t.Fatalf("GetNotes: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d notes, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("note %d = {%d %q, %d-byte body}, want {%d %q, %d-byte body}",
						i, got[i].id, got[i].title, len(got[i].body),
						tt.want[i].id, tt.want[i].title, len(tt.want[i].body))
				}
			}
		})
	}
}

// TestSaveNoteGeneratesID checks that a note with id 0 is stored under a
// fresh non-zero id, and that two such saves do not overwrite each other.
func TestSaveNoteGeneratesID(t *testing.T) {
	s := newTestStore(t)
	for _, title := range []string{"first", "second"} {
		if err := s.SaveNote(Note{title: title, body: "b"}); err != nil {
			t.Fatalf("SaveNote: %v", err)
		}
	}
	got, err := s.GetNotes()
	if err != nil {
		t.Fatalf("GetNotes: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d notes, want 2", len(got))
	}
	for _, n := range got {
		if n.id == 0 {
			t.Errorf("note %q stored with id 0", n.title)
		}
	}
}
