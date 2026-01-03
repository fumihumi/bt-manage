package picker

import (
	"sort"
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

	filteredByState := make([]core.Device, 0, len(devices))
	for _, d := range devices {
		switch title {
		case "Connect":
			if d.Connected {
				continue
			}
		case "Disconnect":
			if !d.Connected {
				continue
			}
		}
		filteredByState = append(filteredByState, d)
	}

	sorted := append([]core.Device(nil), filteredByState...)
	sort.SliceStable(sorted, func(i, j int) bool {
		a := sorted[i]
		b := sorted[j]

		// Connect: prefer disconnected first.
		if title == "Connect" {
			if a.Connected != b.Connected {
				return !a.Connected && b.Connected
			}
		}
		// Disconnect: prefer connected first.
		if title == "Disconnect" {
			if a.Connected != b.Connected {
				return a.Connected && !b.Connected
			}
		}
		// Repair: prefer connected first (paired-device selection).
		if strings.HasPrefix(title, "Repair:") {
			if a.Connected != b.Connected {
				return a.Connected && !b.Connected
			}
		}

		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	m := model{
		title:   title,
		devices: sorted,
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
		case "up", "ctrl+p":
			if m.index > 0 {
				m.index--
			}
			return m, nil
		case "down", "ctrl+n":
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
