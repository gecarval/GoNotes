package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type model struct {
	state     uint
	store     *Store
	notes     []Note
	currNote  Note
	listIndex int
	text      string
}

func newModel(store *Store) model {
	notes, err := store.GetNotes()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return model{
		state: listView,
		store: store,
		notes: notes,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
