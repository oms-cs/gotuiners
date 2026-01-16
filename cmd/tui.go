package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const vPad = 5
const hPad = 10

const (
	ctrs status = iota
	imgs
	details
)

var (
	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99"))

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	tableStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("241"))
)

type Model struct {
	containers table.Model
	images     table.Model
	details    viewport.Model
	focused    status
}

type container struct {
	Name   string `json:"Names"`
	ID     string `json:"ID"`
	Status string `json:"status"`
	State  string `json:"State"`
	Image  string `json:"Image"`
}

type image struct {
	ID         string `json:"ID"`
	Repository string `json:"Repository"`
	Size       string `json:"Size"`
	Tag        string `json:"Tag"`
}

func InitModel(width, height int) Model {
	h := (height - vPad) / 2
	leftWidth := (width - hPad) / 2
	containers := loadContainers(leftWidth, h)
	images := loadImages(leftWidth, h)
	details := loadViewPoint(leftWidth, height-vPad)

	return Model{containers: containers, images: images, focused: ctrs, details: details}
}

func executeShellCommand(commands ...string) string {
	cmd := exec.Command("docker", commands...)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(output)
}

func getContainerList() []container {
	val := executeShellCommand("ps", "--format", "json")
	containers := make([]container, 0, 10)

	for line := range strings.SplitSeq(val, "\n") {
		if line == "" {
			continue
		}
		var ctr container
		if err := json.Unmarshal([]byte(line), &ctr); err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}
		containers = append(containers, ctr)
	}
	return containers
}

func getImageList() []image {
	val := executeShellCommand("images", "--format", "json")
	images := make([]image, 0, 10)

	for line := range strings.SplitSeq(val, "\n") {
		if line == "" {
			continue
		}
		var img image
		if err := json.Unmarshal([]byte(line), &img); err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}
		images = append(images, img)
	}
	return images
}

func getContainerLogs(id string) string {
	val := executeShellCommand("logs", id)
	return val
}

func getImageLogs(id string) string {
	val := executeShellCommand("inspect", id)
	return val
}

func loadImages(totalWidth, height int) table.Model {
	images := getImageList()

	// set columns for image table

	width := totalWidth / 4
	columns := []table.Column{
		{Title: "ID", Width: width},
		{Title: "Repository", Width: width},
		{Title: "Tag", Width: width},
		{Title: "Size", Width: width},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithHeight(height),
	)

	rows := make([]table.Row, 0, len(images))

	for _, image := range images {
		rows = append(rows, table.Row{image.ID, image.Repository, image.Tag, image.Size})
	}

	t.SetRows(rows)
	t.SetStyles(tableStyles())
	return t
}

func loadContainers(totalWidth, height int) table.Model {
	containers := getContainerList()
	// set columns for image table

	width := totalWidth / 5
	columns := []table.Column{
		{Title: "ID", Width: width},
		{Title: "Name", Width: width},
		{Title: "Status", Width: width},
		{Title: "Image", Width: width},
		{Title: "State", Width: width},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithHeight(height),
		table.WithFocused(true),
	)

	t.SetStyles(tableStyles())

	rows := make([]table.Row, 0, len(containers))

	for _, ctr := range containers {
		rows = append(rows, table.Row{ctr.ID, ctr.Name, ctr.Status, ctr.Image, ctr.State})
	}

	t.SetRows(rows)
	return t
}

func loadViewPoint(width, height int) viewport.Model {
	v := viewport.New(width, height-2)
	v.SetContent("Select a container or image")
	return v
}

/* Implementing Model Interface for bubble tea */

func (m Model) Init() tea.Cmd {
	return nil
}

func tableStyles() table.Styles {
	s := table.DefaultStyles()

	s.Header = s.Header.
		Foreground(lipgloss.Color("39")). // Hot Pink text
		Bold(true)                        // Make the Pink headers pop

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("42")).
		Bold(true)

	return s
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = InitModel(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", tea.KeyCtrlC.String():
			return m, tea.Quit

		case "tab":
			m.focused = (m.focused + 1) % 3

			//Set Both Blur
			m.containers.Blur()
			m.images.Blur()

			if m.focused == ctrs {
				m.containers.Focus()
			} else {
				m.images.Focus()
			}

		case "enter":
			switch m.focused {
			case ctrs:
				row := m.containers.SelectedRow()
				if len(row) > 0 {
					m.details.SetContent(getContainerLogs(row[0]))
					m.focused = details
				}
			case imgs:
				row := m.images.SelectedRow()
				if len(row) > 0 {
					m.details.SetContent(getImageLogs(row[0]))
					m.focused = details
				}
			}
			return m, nil
		}
	}

	switch m.focused {
	case ctrs:
		m.containers, cmd = m.containers.Update(msg)
	case imgs:
		m.images, cmd = m.images.Update(msg)
	case details:
		m.details, cmd = m.details.Update(msg)
	}
	cmds = append(cmds, cmd)

	cmd = tea.Batch(cmds...)
	return m, cmd
}

func (m Model) View() string {
	cStyle, iStyle, dStyle := tableStyle, tableStyle, tableStyle

	switch m.focused {
	case ctrs:
		cStyle = focusedStyle
	case imgs:
		iStyle = focusedStyle
	case details:
		dStyle = focusedStyle
	}

	leftElement := lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Left,
			headerStyle.Render(" Containers"),
			cStyle.Render(m.containers.View()),
		),
		lipgloss.JoinVertical(lipgloss.Left,
			headerStyle.Render(" Images"),
			iStyle.Render(m.images.View()),
		),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftElement,
		lipgloss.JoinVertical(lipgloss.Left,
			headerStyle.Render(" Details"),
			dStyle.Render(m.details.View()),
		),
	)
}

func main() {
	p := tea.NewProgram(Model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Failed to Load TUI ", err)
		os.Exit(1)
	}
}
