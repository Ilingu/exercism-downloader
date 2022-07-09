package spinner

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinnerProgram struct {
	Model spinnerModel
}

func (sp SpinnerProgram) SpawnSpinner() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}

type errMsg error

type spinnerModel struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func initialModel() spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return spinnerModel{spinner: s}
}

func (s spinnerModel) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			s.quitting = true
			return s, tea.Quit
		default:
			return s, nil
		}

	case errMsg:
		s.err = msg
		return s, nil

	default:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	}

}

func (s spinnerModel) View() string {
	if s.err != nil {
		return s.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Downloading exercism...\n\n", s.spinner.View())
	if s.quitting {
		return str + "\n"
	}
	return str
}
