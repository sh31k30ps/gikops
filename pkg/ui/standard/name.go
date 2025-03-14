package standard

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func PromptName(placeholder string, forWhat string) (string, error) {
	p := tea.NewProgram(newPromptModel(placeholder, fmt.Sprintf("What is the name for the %s?", forWhat)))
	m, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running name selection: %w", err)
	}
	nameModel, ok := m.(promptModel)
	if !ok {
		return "", fmt.Errorf("could not convert name model")
	}
	if nameModel.textInput.Value() != "" {
		return nameModel.textInput.Value(), nil
	}
	return placeholder, nil
}
