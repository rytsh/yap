package login

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/rytsh/yap/internal/tui/style"
)

var (
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(style.Highlight).
			BorderTop(false).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	bannerStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(false).
			BorderRight(false).
			Padding(0, 1)
)
