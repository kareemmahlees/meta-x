package utils

import "github.com/charmbracelet/lipgloss"

func NewStyle(text, bgColorHex string) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(bgColorHex)).
		MarginTop(1).
		SetString(text)

}
