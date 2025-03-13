package standard

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func PromptName(placeholder string, forWhat string) (string, error) {
	p := tea.NewProgram(newNameModel(placeholder, forWhat))
	m, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running name selection: %w", err)
	}
	nameModel, ok := m.(nameModel)
	if !ok {
		return "", fmt.Errorf("could not convert name model")
	}
	if nameModel.textInput.Value() != "" {
		return nameModel.textInput.Value(), nil
	}
	return placeholder, nil
}

type (
	errMsg error
)

type nameModel struct {
	textInput textinput.Model
	err       error
	forWhat   string
}

func newNameModel(defaultName string, forWhat string) nameModel {
	ti := textinput.New()
	ti.Placeholder = defaultName
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return nameModel{
		textInput: ti,
		err:       nil,
		forWhat:   forWhat,
	}
}

func (m nameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m nameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m nameModel) View() string {
	return fmt.Sprintf(
		"What is the name for the %s?\n\n%s",
		m.forWhat,
		m.textInput.View(),
	)
}
