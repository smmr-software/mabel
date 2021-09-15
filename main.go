package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func (m model) Init() tea.Cmd {
	return tick()
}

func main() {
	model, err := initialModel()
	if err != nil {
		log.Fatal(err)
	}

	if err := tea.NewProgram(model, tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
