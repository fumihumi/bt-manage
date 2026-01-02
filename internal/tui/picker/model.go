package picker

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fumihumi/bt-manage/internal/core"
)

type model struct {
	title   string
	devices []core.Device

	input textinput.Model

	filtered []core.Device
	index    int

	selected core.Device
	canceled bool
	chosen   bool

	width  int
	height int
}

func newModel(title string, devices []core.Device) model {
	in := textinput.New()
	in.Placeholder = "search"
	in.Focus()
	in.CharLimit = 256
	in.Width = 30

	m := model{
		title:   title,
		devices: append([]core.Device(nil), devices...),
		input:   in,
		index:   0,
	}
	m.applyFilter()
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = min(50, max(20, msg.Width-10))
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.canceled = true
			m.chosen = false
			return m, tea.Quit
		case "enter":
			if len(m.filtered) == 0 {
				return m, nil
			}
			m.selected = m.filtered[m.index]
			m.chosen = true
			return m, tea.Quit
		case "up", "k":
			if m.index > 0 {
				m.index--
			}
			return m, nil
		case "down", "j":
			if m.index < len(m.filtered)-1 {
				m.index++
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.applyFilter()
	if m.index >= len(m.filtered) {
		m.index = max(0, len(m.filtered)-1)
	}
	return m, cmd
}

func (m *model) applyFilter() {
	q := strings.TrimSpace(strings.ToLower(m.input.Value()))
	if q == "" {
		m.filtered = append([]core.Device(nil), m.devices...)
		return
	}

	out := make([]core.Device, 0, len(m.devices))
	for _, d := range m.devices {
		if strings.Contains(deviceSearchKey(d), q) {
			out = append(out, d)
		}
	}
	m.filtered = out
	if m.index >= len(m.filtered) {
		m.index = max(0, len(m.filtered)-1)
	}
}

func deviceSearchKey(d core.Device) string {
	// Lowercase, concatenated for simple substring match.
	// Address may be empty depending on backend; keep it safe.
	return strings.ToLower(strings.TrimSpace(d.Name + " " + d.Address))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
