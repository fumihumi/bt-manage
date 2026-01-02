package picker

import (
	"sort"
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
	// Avoid textinput treating Enter as submit; we want Enter to confirm selection.
	// (bubbles/textinput has no EnterKeySubmits; keep behavior by handling Enter in Update.)

	// Prevent textinput's suggestion navigation from stealing up/down.
	in.KeyMap.NextSuggestion.SetKeys()
	in.KeyMap.PrevSuggestion.SetKeys()

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

		if title == "Connect" {
			if a.Connected != b.Connected {
				return !a.Connected && b.Connected
			}
		}
		if title == "Disconnect" {
			if a.Connected != b.Connected {
				return a.Connected && !b.Connected
			}
		}

		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	m := multiModel{
		title:       title,
		devices:     sorted,
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
	// First handle key events we must own, before delegating to textinput.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = min(50, max(20, msg.Width-10))
		return m, nil
	case tea.KeyMsg:
		// KeyType-based handling.
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.canceled = true
			m.chosen = false
			return m, tea.Quit
		case tea.KeyEnter:
			m.selected = m.selected[:0]
			for _, d := range m.filtered {
				if m.selectedMap[d.Address] {
					m.selected = append(m.selected, d)
				}
			}
			m.chosen = true
			return m, tea.Quit
		case tea.KeyUp:
			if m.index > 0 {
				m.index--
			}
			return m, nil
		case tea.KeyDown:
			if m.index < len(m.filtered)-1 {
				m.index++
			}
			return m, nil
		case tea.KeySpace:
			if len(m.filtered) == 0 {
				return m, nil
			}
			d := m.filtered[m.index]
			m.selectedMap[d.Address] = !m.selectedMap[d.Address]
			return m, nil
		}

		// Some terminals send space/j as runes; keep a string fallback.
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
