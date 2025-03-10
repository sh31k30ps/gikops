package setup

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type nameModel struct {
	textInput textinput.Model
	err       error
}

func newNameModel(defaultName string) nameModel {
	ti := textinput.New()
	ti.Placeholder = defaultName
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return nameModel{
		textInput: ti,
		err:       nil,
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
		"What is the name of your project?\n\n%s",
		m.textInput.View(),
	)
}
