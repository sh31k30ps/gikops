package standard

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

func Prompt(placeholder, label string) (string, error) {
	p := tea.NewProgram(newPromptModel(placeholder, label))
	m, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running name selection: %w", err)
	}
	promptModel, ok := m.(promptModel)
	if !ok {
		return "", fmt.Errorf("could not convert name model")
	}
	if promptModel.textInput.Value() != "" {
		return promptModel.textInput.Value(), nil
	}
	return placeholder, nil
}

type promptModel struct {
	textInput textinput.Model
	err       error
	label     string
}

func newPromptModel(defaultName string, label string) promptModel {
	ti := textinput.New()
	ti.Placeholder = defaultName
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return promptModel{
		textInput: ti,
		err:       nil,
		label:     label,
	}
}

func (m promptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m promptModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.label,
		m.textInput.View(),
	)
}
