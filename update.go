package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		var cmd tea.Cmd
		cmds := make([]tea.Cmd, 5)
		cmd = modelStateMachine(&m, msg.String())
		cmds = append(cmds, cmd)
		m.textinput, cmd = m.textinput.Update(msg)
		cmds = append(cmds, cmd)
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func modelStateMachine(m *model, key string) tea.Cmd {
	switch m.state {
	case listView:
		return listViewUpdate(m, key)
	case titleView:
		return titleViewUpdate(m, key)
	case bodyView:
		return bodyViewUpdate(m, key)
	}
	return nil
}

func listViewUpdate(m *model, key string) tea.Cmd {
	switch key {
	case "ctrl+c":
		return tea.Quit
	case "n":
		m.state = titleView
		m.currNote = Note{}
		m.textinput.SetValue("")
		m.textarea.SetValue("")
		m.textinput.Focus()
	}
	return nil
}

func titleViewUpdate(m *model, key string) tea.Cmd {
	switch key {
	case "ctrl+c":
		m.textarea.Blur()
		m.textinput.Blur()
		m.state = listView
	case "enter":
		if m.textinput.Value() != "" {
			m.state = bodyView
			m.currNote.title = m.textinput.Value()
			m.textinput.Blur()
			m.textarea.SetValue("")
			m.textarea.Focus()
		}
	}
	return nil
}

func bodyViewUpdate(m *model, key string) tea.Cmd {
	switch key {
	case "ctrl+c":
		m.state = titleView
		m.textarea.Blur()
		m.textinput.Focus()
	case "ctrl+s":
		m.currNote.body = m.textarea.Value()
		if err := m.store.SaveNote(&m.currNote); err != nil {
			log.Fatalf("Error: %v", err)
		}
	case "ctrl+z":
		m.textarea.SetValue(m.currNote.body)
	default:
		m.textarea.Value()
	}
	return nil
}
