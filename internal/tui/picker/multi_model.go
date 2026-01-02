package picker

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fumihumi/bt-manage/internal/core"
)

type multiModel struct {
	title   string
	devices []core.Device

	input textinput.Model

	filtered []core.Device
	index    int

	selectedMap map[string]bool
	selected    []core.Device
	canceled    bool
	chosen      bool

	width  int
	height int
}

func newMultiModel(title string, devices []core.Device) multiModel {
	in := textinput.New()
	in.Placeholder = "search"
	in.Focus()
	in.CharLimit = 256
	in.Width = 30

	m := multiModel{
		title:       title,
		devices:     append([]core.Device(nil), devices...),
		input:       in,
		index:       0,
		selectedMap: map[string]bool{},
	}
	m.applyFilter()
	return m
}

func (m multiModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m multiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.selected = m.selected[:0]
			for _, d := range m.filtered {
				if m.selectedMap[d.Address] {
					m.selected = append(m.selected, d)
				}
			}
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
		case " ":
			if len(m.filtered) == 0 {
				return m, nil
			}
			d := m.filtered[m.index]
			// Address is used as a stable key.
			m.selectedMap[d.Address] = !m.selectedMap[d.Address]
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

func (m *multiModel) applyFilter() {
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
