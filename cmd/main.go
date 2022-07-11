package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"exercism-cli/cmd/subs/spinner"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var mainProgram *tea.Program

func main() {
	InitialName := flag.String("n", "", "The exercism name")
	flag.Parse()

	mainProgram = tea.NewProgram(initMainModel(*InitialName))
	if err := mainProgram.Start(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}

type ClearErrorMsg bool
type errMsg error

// The main TUI model
type mainModel struct {
	exercismName textinput.Model
	quit         bool
	err          error
}

func initMainModel(InitialName string) mainModel {
	// Input
	ti := textinput.New()
	ti.Placeholder = "Crypto Square"
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#635111"))
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9ca24"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#874BFD")).Bold(true)
	ti.Width = 50

	ti.Focus()
	ti.SetValue(InitialName)

	return mainModel{exercismName: ti}
}

func clearErrTick() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return ClearErrorMsg(true)
	})
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, clearErrTick())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			ExercismName, err := parseInput(m.exercismName.Value())
			if err != nil {
				m.err = err
				return m, nil
			}

			m.exercismName.Blur()
			m.quit = true

			go (spinner.SpinnerProgram{}).SpawnSpinner() // Render Spinner
			go DownloadAndOpenExcerism(ExercismName)

			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case ClearErrorMsg:
		m.err = nil
		return m, clearErrTick()

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.exercismName, cmd = m.exercismName.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	if m.quit {
		return ""
	}

	var HeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6366f1")).
		PaddingRight(2).PaddingLeft(2).Underline(true)

	var EscapeStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"})

	HeaderUI := HeaderStyle.Render("Enter The Exercism Name:")
	InputUI := fmt.Sprintf(
		"\n\n%s",
		m.exercismName.View(),
	)
	FooterUI := fmt.Sprintf("\n\n%s\n", EscapeStyle.Render("(esc to quit)"))

	if m.err != nil {
		ErrorUI := fmt.Sprintf("\n\n%s", lipgloss.NewStyle().Foreground(lipgloss.Color("#eb4d4b")).Render(m.err.Error()))
		return HeaderUI + InputUI + ErrorUI + FooterUI
	}
	return HeaderUI + InputUI + FooterUI
}
