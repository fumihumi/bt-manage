package picker

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("type to filter • ↑/↓ (j/k) move • enter select • esc cancel"))
	b.WriteString("\n\n")

	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	if len(m.filtered) == 0 {
		b.WriteString(dimStyle.Render("(no matches)"))
		b.WriteString("\n")
		return b.String()
	}

	// show up to N items
	maxItems := max(5, min(15, m.height-8))
	start := 0
	if m.index >= maxItems {
		start = m.index - maxItems + 1
	}
	end := min(len(m.filtered), start+maxItems)

	for i := start; i < end; i++ {
		prefix := "  "
		line := fmt.Sprintf("%s", m.filtered[i].Name)
		if i == m.index {
			prefix = "> "
			line = selectedStyle.Render(line)
		}
		b.WriteString(prefix)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}
