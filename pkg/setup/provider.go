package setup

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Model for provider selection
type providerModel struct {
	choices  []string
	cursor   int
	selected int
	quitting bool
}

func newProviderModel(providers []string) providerModel {
	return providerModel{
		choices:  providers,
		selected: 0, // Select first by default
	}
}

func (m providerModel) Init() tea.Cmd {
	return nil
}

func (m providerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.cursor
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m providerModel) View() string {
	s := "Select container runtime provider:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}
