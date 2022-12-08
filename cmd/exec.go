package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Recommanded to call this function in a go routines, it will downloads the exercism to the user local system (if exercism cli is configured) and will then open it in vscodium
func DownloadAndOpenExcerism(exercismName string) {
	SendFinalResponse := func(res string, isErr bool) {
		go PushMessage(res, isErr)
		time.Sleep(2 * time.Second)
		mainProgram.Quit() // Quit main program
	}

	dlPath, err := DownloadExercism(exercismName)
	if err != nil {
		SendFinalResponse("failed to download exercism", true)
		return
	}

	var try uint
	for {
		if try >= 5 {
			SendFinalResponse("failed to download exercism", true)
			return
		}

		if _, err := os.Stat(dlPath); err == nil {
			break
		}

		try++
		go PushMessage("Folder not created yet, backing up 1s...", true)
		time.Sleep(1 * time.Second)
	}
	go PushMessage(fmt.Sprintf("%s downloaded, opening it...", exercismName), false)

	err = OpenEditor(dlPath)
	if err != nil {
		SendFinalResponse("failed to open editor", true)
		return
	}

	SendFinalResponse("operation terminate successfully", false)
}

func DownloadExercism(exercismName string) (string, error) {
	ExercismArgs := []string{"download", fmt.Sprintf("--exercise=%s", exercismName), "--track=go"}
	dlPath, err := ExecCmdWithOutput("exercism", ExercismArgs...) // Download exercism via exercism cli
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(dlPath), nil
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
