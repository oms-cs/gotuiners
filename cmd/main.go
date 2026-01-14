package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	table         table.Model
	allContainers []Container
	loading       bool
	err           error
}

type Container struct {
	Name   string `json:"Names"`
	ID     string `json:"ID"`
	Image  string `json:"Image"`
	Status string `json:"Status"`
	State  string `json:"State"`
}

var baseStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("99")).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("29"))

func (m model) Init() tea.Cmd {
	return getContainers
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ContainerMsg:
		rows := make([]table.Row, 0, len(msg))

		for _, ctr := range msg {
			rows = append(rows, table.Row{ctr.Name, ctr.ID, ctr.Image, ctr.Status, ctr.State})
		}
		m.allContainers = msg
		m.table.SetRows(rows)
	case tea.KeyMsg:
		switch msg.Type.String() {
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			return m, tea.Quit
		case "enter":
			index := m.table.Cursor()
			selectedContainer := m.allContainers[index]
			// Now you have the full struct with all hidden fields
			return m, tea.Batch(
				tea.Printf("Selected Name: %s, Image: %s", selectedContainer.Name, selectedContainer.Image),
			)
		}
		if len(msg.Runes) > 0 && msg.Runes[0] == 'q' {
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return "\n  Docker Containers\n\n" + baseStyle.Render(m.table.View()) + "\n  press q to quit\n"
}

func getContainers() tea.Msg {
	jsonString, err := executeCommand()
	if err != nil {
		return errMsg{err}
	}
	lines := strings.Split(strings.TrimSpace(string(jsonString)), "\n")
	var containers []Container

	for _, line := range lines {
		if line == "" {
			continue
		}
		var c Container
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}
		containers = append(containers, c)
	}

	return ContainerMsg(containers)
}

func executeCommand() (string, error) {

	cmd := exec.Command("docker", "ps", "--format", "json")
	res, err := cmd.Output()

	if err != nil {
		// There was an error making our request. Wrap the error we received
		// in a message and return it.
		return "", err
	}
	// We received a response from the server. Return the HTTP status code
	// as a message.
	return string(res), nil
}

type errMsg struct{ err error }

type ContainerMsg []Container

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func getTable() table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "ID", Width: 20},
		{Title: "Image", Width: 20},
		{Title: "Status", Width: 20},
		{Title: "State", Width: 20},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	// Apply some Bubblegum styling
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Bold(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("240")).Bold(false)
	t.SetStyles(s)
	return t
}

func main() {
	t := getTable()
	p := tea.NewProgram(model{table: t})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
