package main

import (
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
			got, cmd := newModel().Update(tt.msg)
			if gotQuit := cmd != nil && isQuit(cmd()); gotQuit != tt.wantQuit {
				t.Fatalf("quit = %v, want %v", gotQuit, tt.wantQuit)
			}
			if got.(model) != newModel() {
				t.Fatalf("model changed: %+v", got)
			}
		})
	}
}

// TestUpdateIgnoresUnknownMsg makes sure non-key messages are a no-op.
func TestUpdateIgnoresUnknownMsg(t *testing.T) {
	if _, cmd := newModel().Update(tea.WindowSizeMsg{Width: 1, Height: 1}); cmd != nil {
		t.Fatal("expected nil cmd for window size message")
	}
}

// TestViewEmpty checks that a fresh model renders without panicking.
func TestViewEmpty(t *testing.T) {
	_ = newModel().View()
}

func isQuit(msg tea.Msg) bool {
	_, ok := msg.(tea.QuitMsg)
	return ok
}
