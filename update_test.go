package main

import (
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// TestUpdateKeys checks that only ctrl+c quits and that other keys leave
// the model untouched.
func TestUpdateKeys(t *testing.T) {
	tests := []struct {
		name     string
		msg      tea.KeyPressMsg
		wantQuit bool
	}{
		{"ctrl+c quits", tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}, true},
		{"plain c ignored", tea.KeyPressMsg{Code: 'c', Text: "c"}, false},
		{"q ignored", tea.KeyPressMsg{Code: 'q', Text: "q"}, false},
		{"enter ignored", tea.KeyPressMsg{Code: tea.KeyEnter}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStore(t)
			got, cmd := newModel(s).Update(tt.msg)
			if gotQuit := cmd != nil && isQuit(cmd()); gotQuit != tt.wantQuit {
				t.Fatalf("quit = %v, want %v", gotQuit, tt.wantQuit)
			}
			// model holds a slice, so it is not comparable with ==.
			if want := newModel(s); !reflect.DeepEqual(got.(model), want) {
				t.Fatalf("model changed: %+v", got)
			}
		})
	}
}

// TestUpdateIgnoresUnknownMsg makes sure non-key messages are a no-op.
func TestUpdateIgnoresUnknownMsg(t *testing.T) {
	m := newModel(newTestStore(t))
	if _, cmd := m.Update(tea.WindowSizeMsg{Width: 1, Height: 1}); cmd != nil {
		t.Fatal("expected nil cmd for window size message")
	}
}

// TestNewModel checks that the model starts in the list view, holds the
// store it was given, and loads the notes already saved in it.
func TestNewModel(t *testing.T) {
	tests := []struct {
		name  string
		saved []Note
	}{
		{"empty store", nil},
		{"one note", []Note{{1, "t", "b"}}},
		{"several notes", []Note{{1, "a", "x"}, {2, "日本語", "🐐"}, {3, "", ""}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStore(t)
			for _, n := range tt.saved {
				if err := s.SaveNote(n); err != nil {
					t.Fatalf("SaveNote: %v", err)
				}
			}
			m := newModel(s)
			if m.store != s {
				t.Fatal("model does not hold the given store")
			}
			if m.state != listView {
				t.Fatalf("state = %d, want listView", m.state)
			}
			if m.listIndex != 0 {
				t.Fatalf("listIndex = %d, want 0", m.listIndex)
			}
			if !reflect.DeepEqual(m.notes, tt.saved) {
				t.Fatalf("notes = %+v, want %+v", m.notes, tt.saved)
			}
		})
	}
}

// TestViewEmpty checks that a fresh model renders without panicking.
func TestViewEmpty(t *testing.T) {
	_ = newModel(newTestStore(t)).View()
}

func isQuit(msg tea.Msg) bool {
	_, ok := msg.(tea.QuitMsg)
	return ok
}
