package main

import (
	tea "charm.land/bubbletea/v2"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type model struct {
	state uint
	text  string
}

func newModel() model {
	return model{
		state: listView,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
