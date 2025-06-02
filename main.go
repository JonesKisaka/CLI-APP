package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type model struct {
	table table.Model
	err   error
}

func getRunningContainers() ([]table.Row, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: false})
	if err != nil {
		return nil, err
	}

	var rows []table.Row
	for _, c := range containers {
		ports := ""
		for _, p := range c.Ports {
			ports += fmt.Sprintf("%s:%d->%d/%s ", p.IP, p.PublicPort, p.PrivatePort, p.Type)
		}
		name := c.Names[0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}

		rows = append(rows, table.Row{
			c.ID[:12],
			name,
			c.Image,
			c.Status,
			ports,
		})
	}

	return rows, nil
}

func initialModel() model {
	columns := []table.Column{
		{Title: "ID", Width: 12},
		{Title: "Name", Width: 20},
		{Title: "Image", Width: 25},
		{Title: "Status", Width: 20},
		{Title: "Ports", Width: 30},
	}

	rows, err := getRunningContainers()
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	t.SetStyles(table.DefaultStyles())

	return model{
		table: t,
		err:   err,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	return "\nRunning Docker Containers:\n\n" + m.table.View() + "\n\nPress q to quit."
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting program: %v", err)
		os.Exit(1)
	}
}
