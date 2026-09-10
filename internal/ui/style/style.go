package style

import "github.com/charmbracelet/lipgloss"

var (
	Brand   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	Header  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")).Underline(true)
	Value   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	Dim     = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	Success = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	Warning = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	Error   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	Info    = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	Accent  = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
	Section = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
)