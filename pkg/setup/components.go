package setup

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Model for component selection
type componentModel struct {
	name       string
	choices    []string
	cursor     int
	selected   map[int]bool
	quitting   bool
	components []string
}

func newComponentModel(name string, components []string) componentModel {
	selected := make(map[int]bool)
	// Select all by default
	for i := range components {
		selected[i] = true
	}
	return componentModel{
		name:       name,
		choices:    components,
		selected:   selected,
		components: components,
	}
}

func (m componentModel) Init() tea.Cmd {
	return nil
}

func (m componentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m componentModel) View() string {
	s := fmt.Sprintf("Select \033[1m%s\033[0m components to install (space to toggle, enter to confirm):\n\n", m.name)

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.selected[i] {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// s += "\nPress q to quit.\n"

	return s
}
