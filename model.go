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
	store *Store
	text  string
}

func newModel(store *Store) model {
	return model{
		state: listView,
		store: store,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
