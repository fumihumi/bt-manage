package picker

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fumihumi/bt-manage/internal/core"
)

func (m multiModel) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("type to filter (name/address) • ↑/↓ (ctrl+p/ctrl+n) move • space toggle • enter confirm • esc cancel"))
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

	maxItems := max(5, min(15, m.height-8))
	start := 0
	if m.index >= maxItems {
		start = m.index - maxItems + 1
	}
	end := min(len(m.filtered), start+maxItems)

	for i := start; i < end; i++ {
		d := m.filtered[i]
		checked := "[ ]"
		if m.selectedMap[d.Address] {
			checked = "[x]"
		}

		prefix := "  "
		nameLine := d.Name
		if nameLine == "" {
			nameLine = "(unknown)"
		}
		meta := deviceMetaMulti(d)

		if i == m.index {
			prefix = "> "
			nameLine = selectedStyle.Render(nameLine)
			checked = selectedStyle.Render(checked)
			meta = selectedStyle.Render(meta)
		} else {
			checked = metaStyle.Render(checked)
			meta = metaStyle.Render(meta)
		}

		b.WriteString(prefix)
		b.WriteString(checked)
		b.WriteString(" ")
		b.WriteString(nameLine)
		b.WriteString("\n")
		if meta != "" {
			b.WriteString("    ")
			b.WriteString(meta)
			b.WriteString("\n")
		}
	}

	return b.String()
}

func deviceMetaMulti(d core.Device) string {
	// reuse single picker meta format
	parts := make([]string, 0, 2)
	if strings.TrimSpace(d.Address) != "" {
		parts = append(parts, d.Address)
	}
	if d.Connected {
		parts = append(parts, "connected")
	}
	return strings.Join(parts, " • ")
}
