package projectcomponent

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	questionSelectComponents = "Select \033[1m%s\033[0m components to install: \n"
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

func PromptComponents(name string, components []string) ([]string, error) {
	p := tea.NewProgram(newComponentModel(name, components))
	m, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("error running component %s selection: %w", name, err)
	}

	componentModel, ok := m.(componentModel)
	if !ok {
		return nil, fmt.Errorf("could not convert component model")
	}

	if componentModel.quitting {
		return nil, fmt.Errorf("component selection cancelled")
	}

	// Convert selected map to slice
	var selectedComponents []string
	for i, choice := range componentModel.choices {
		if componentModel.selected[i] {
			selectedComponents = append(selectedComponents, choice)
		}
	}
	return selectedComponents, nil
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
	s := fmt.Sprintf(questionSelectComponents, m.name)

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
