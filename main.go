package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func (m model) Init() tea.Cmd {
	return tick()
}

func main() {
	if err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start(); err != nil {
		os.Exit(1)
	}
}
