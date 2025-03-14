package standard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type yesNoModel struct {
	textInput textinput.Model
	question  string
}

func PromptYesNo(question string) (bool, error) {
	p := tea.NewProgram(newYesNoModel(question))
	m, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("error running yes/no prompt: %w", err)
	}

	model, ok := m.(yesNoModel)
	if !ok {
		return false, fmt.Errorf("could not convert yes/no model")
	}

	answer := strings.ToLower(strings.TrimSpace(model.textInput.Value()))
	return answer == "y" || answer == "yes", nil
}

func newYesNoModel(question string) yesNoModel {
	ti := textinput.New()
	ti.Placeholder = "y/N"
	ti.Focus()
	ti.CharLimit = 3
	ti.Width = 3

	return yesNoModel{
		textInput: ti,
		question:  question,
	}
}

func (m yesNoModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m yesNoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m yesNoModel) View() string {
	return fmt.Sprintf(
		"%s [y/N]\n\n%s",
		m.question,
		m.textInput.View(),
	)
}
