package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	subresult "exercism-cli/cmd/subs/result"

	tea "github.com/charmbracelet/bubbletea"
)

// Recommanded to call this function in a go routines, it will downloads the exercism to the user local system (if exercism cli is configured) and will then open it in vscodium
func DownloadAndOpenExcerism(exercismName string) {
	SendResponse := func(res string, isErr bool) {
		NewResponse := subresult.ResultSubProgram{Model: subresult.ResultModel{Result: res, IsError: isErr}}

		go NewResponse.SpawnResultSubProgram() // Render result to ui
		time.Sleep(2 * time.Second)

		mainProgram.Quit() // Quit main program
	}

	dlPath, err := DownloadExercism(exercismName)
	if err != nil {
		SendResponse("failed to download exercism", true)
		return
	}

	err = OpenEditor(dlPath)
	if err != nil {
		SendResponse("failed to open editor", true)
		return
	}

	SendResponse("operation terminate successfully", false)
}

func DownloadExercism(exercismName string) (string, error) {
	ExercismArgs := []string{"download", fmt.Sprintf("--exercise=%s", exercismName), "--track=go"}
	dlPath, err := ExecCmdWithOutput("exercism", ExercismArgs...) // Download exercism via exercism cli
	if err != nil {
		return "", err
	}

	return dlPath, nil
}

func OpenEditor(dlPath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" && strings.Contains(dlPath, "Iling") {
		editor = "codium"
	} else if editor == "" {
		editor = "code"
	}

	err := ExecCmd(editor, dlPath) // Open editor
	if err != nil {
		return err
	}

	return nil
}

// Helpers
type FailedCmdMsg struct{ err error }

func CreateTeaCmd(name string, args ...string) tea.Cmd {
	c := exec.Command(name, args...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		fmt.Println(err)
		return FailedCmdMsg{err}
	})
}

func ExecCmdWithOutput(name string, args ...string) (string, error) {
	c := exec.Command(name, args...)

	data, err := c.Output()
	if err != nil {
		return "", errors.New("cannot exec cmd")
	}
	return string(data), nil
}

func ExecCmd(name string, args ...string) error {
	c := exec.Command(name, args...)

	err := c.Start()
	if err != nil {
		return errors.New("cannot start cmd")
	}

	err = c.Wait()
	if err != nil {
		return errors.New("cannot wait cmd")
	}
	return nil
}
