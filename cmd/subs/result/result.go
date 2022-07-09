package subresult

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResultSubProgram struct {
	Model ResultModel
}

type ResultModel struct {
	Result  string
	IsError bool
}

var program *tea.Program

// It spawns a new tea TUI and display the result
func (sp ResultSubProgram) SpawnResultSubProgram() {
	program = tea.NewProgram(sp.Model)
	if err := program.Start(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}

func (r ResultModel) Init() tea.Cmd {
	return nil
}

func (r ResultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyEnter:
			return r, tea.Quit
		}
	}
	return r, nil
}

func (r ResultModel) View() string {
	var baseStyle = lipgloss.NewStyle().PaddingRight(2).PaddingLeft(2).Underline(true).Bold(true)

	var SuccessStyle = baseStyle.Copy().Foreground(lipgloss.Color("#43BF6D"))
	var ErrorStyle = baseStyle.Copy().Foreground(lipgloss.Color("#eb4d4b"))

	if r.IsError {
		return ErrorStyle.Render(r.Result + " ❌")
	}
	return SuccessStyle.Render(r.Result + " ✅")
}
