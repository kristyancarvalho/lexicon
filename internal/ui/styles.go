package ui

import "github.com/charmbracelet/lipgloss"

const (
	ansiRed     = lipgloss.ANSIColor(1)
	ansiGreen   = lipgloss.ANSIColor(2)
	ansiYellow  = lipgloss.ANSIColor(3)
	ansiBlue    = lipgloss.ANSIColor(4)
	ansiMagenta = lipgloss.ANSIColor(5)
	ansiCyan    = lipgloss.ANSIColor(6)
	ansiGray    = lipgloss.ANSIColor(8)
)

var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(ansiMagenta)
	mutedStyle     = lipgloss.NewStyle().Foreground(ansiGray)
	userStyle      = lipgloss.NewStyle().Bold(true).Foreground(ansiCyan)
	assistantStyle = lipgloss.NewStyle().Bold(true).Foreground(ansiYellow)
	errorStyle     = lipgloss.NewStyle().Foreground(ansiRed)
	okStyle        = lipgloss.NewStyle().Foreground(ansiGreen)
	activeStyle    = lipgloss.NewStyle().Bold(true).Foreground(ansiMagenta)
	accentStyle    = lipgloss.NewStyle().Foreground(ansiBlue)
	ruleStyle      = lipgloss.NewStyle().Foreground(ansiGray)
	keyStyle       = lipgloss.NewStyle().Bold(true)
	userBarStyle   = lipgloss.NewStyle().Border(lipgloss.ThickBorder(), false, false, false, true).BorderForeground(ansiCyan).PaddingLeft(1)
	botBarStyle    = lipgloss.NewStyle().Border(lipgloss.ThickBorder(), false, false, false, true).BorderForeground(ansiYellow).PaddingLeft(1)
	modalStyle     = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ansiMagenta).
			Padding(1, 3)
)
