package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	var content string

	if m.state == listView {
		newNoteStringFormat := "Notes List (press 'n' for new note)"
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			newNoteStringFormat,
		)
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			"Title:",
			m.textinput.View(),
			"",
			"Body (ctrl+s to save, esc to cancel):",
			m.textarea.View(),
		)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
