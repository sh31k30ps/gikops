package standard

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Model for provider selection
type ChoiceModel struct {
	label    string
	choices  []string
	cursor   int
	selected int
	quitting bool
}

func PromptChoice(label string, choices []string) (string, error) {
	p := tea.NewProgram(newChoiceModel(label, choices))
	m, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running choice prompt: %w", err)
	}
	choiceModel, ok := m.(ChoiceModel)
	if !ok {
		return "", fmt.Errorf("could not convert choice model")
	}
	if choiceModel.quitting {
		return "", fmt.Errorf("choice prompt cancelled")
	}
	return choiceModel.choices[choiceModel.selected], nil
}

func newChoiceModel(label string, choices []string) ChoiceModel {
	return ChoiceModel{
		label:    label,
		choices:  choices,
		selected: 0, // Select first by default
	}
}

func (m ChoiceModel) Init() tea.Cmd {
	return nil
}

func (m ChoiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m ChoiceModel) View() string {
	s := fmt.Sprintf("%s:\n\n", m.label)

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
