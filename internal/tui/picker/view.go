package picker

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fumihumi/bt-manage/internal/core"
)

func (m model) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("type to filter (name/address) • ↑/↓ (ctrl+p/ctrl+n) move • enter select • esc cancel"))
	b.WriteString("\n\n")

	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	if len(m.filtered) == 0 {
		q := strings.TrimSpace(m.input.Value())
		if q == "" {
			b.WriteString(dimStyle.Render("(no devices found)"))
		} else {
			b.WriteString(dimStyle.Render(fmt.Sprintf("(no matches for %q)", q)))
		}
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
		d := m.filtered[i]
		prefix := "  "
		nameLine := d.Name
		if nameLine == "" {
			nameLine = "(unknown)"
		}
		meta := deviceMeta(d)

		if i == m.index {
			prefix = "> "
			nameLine = selectedStyle.Render(nameLine)
			meta = selectedStyle.Render(meta)
		} else {
			meta = metaStyle.Render(meta)
		}

		b.WriteString(prefix)
		b.WriteString(nameLine)
		if meta != "" {
			b.WriteString("\n   ")
			b.WriteString(meta)
		}
		b.WriteString("\n")
	}

	return b.String()
}

func deviceMeta(d core.Device) string {
	parts := make([]string, 0, 2)
	if strings.TrimSpace(d.Address) != "" {
		parts = append(parts, d.Address)
	}
	if d.Connected {
		parts = append(parts, "connected")
	}
	return strings.Join(parts, " • ")
}
