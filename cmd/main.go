package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	response string
	err      error
}

func (m model) Init() tea.Cmd {
	return executeCommand
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		m.response = string(msg)
		return m, tea.Quit

	case errMsg:
		// There was an error. Note it in the model. And tell the runtime
		// we're done and want to quit.
		m.err = msg
		return m, tea.Quit

	// check if any key was pressed
	case tea.KeyMsg:
		//if it was key press then do key press actions here
		switch msg.String() {
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	// Tell the user we're doing something.
	s := "Loading ...\n"

	// When the server responds with a status, add it to the current line.
	if len(m.response) > 0 {
		s += fmt.Sprintf("%s", m.response)
	}

	// Send off whatever we came up with above for rendering.
	return "\n" + s + "\n\n"
}

func executeCommand() tea.Msg {

	cmd := exec.Command("docker", "ps")
	res, err := cmd.Output()

	if err != nil {
		// There was an error making our request. Wrap the error we received
		// in a message and return it.
		return errMsg{err}
	}
	// We received a response from the server. Return the HTTP status code
	// as a message.
	return statusMsg(string(res))
}

type statusMsg string

type errMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
