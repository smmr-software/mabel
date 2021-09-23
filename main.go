package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
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
