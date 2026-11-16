package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	mutedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	userStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	assistantStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("221"))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	okStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
	activeStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	modalStyle     = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("212")).
			Padding(1, 3)
)
