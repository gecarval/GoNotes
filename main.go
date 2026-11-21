package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

func main() {
	s := &Store{}
	if err := s.Init(); err != nil {
		log.Fatalf("Error: %v", err)
	}
	m := newModel(s)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
