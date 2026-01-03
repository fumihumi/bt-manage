package picker

import (
	"context"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fumihumi/bt-manage/internal/core"
)

type devicesUpdateMsg struct {
	devices []core.Device
}

type spinnerTickMsg struct{}

type streamModel struct {
	model
	spinning bool
	dots     int
}

func newStreamModel(title string) streamModel {
	m := streamModel{model: newModel(title, nil), spinning: true, dots: 0}
	// Keep placeholder search focused.
	m.input.Focus()
	return m
}

func (m streamModel) Init() tea.Cmd {
	return tea.Batch(textinputBlink(), spinnerTick())
}

func spinnerTick() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(time.Time) tea.Msg { return spinnerTickMsg{} })
}

func textinputBlink() tea.Cmd {
	return func() tea.Msg { return nil }
}

func (m streamModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case devicesUpdateMsg:
		// Normalize by address and replace the list to avoid UI duplication/flicker.
		m.devices = normalizeDevices(msg.devices)
		m.applyFilter()
		if len(m.devices) > 0 {
			m.spinning = false
		}
		if m.index >= len(m.filtered) {
			m.index = max(0, len(m.filtered)-1)
		}
		return m, nil
	case spinnerTickMsg:
		if m.spinning {
			m.dots = (m.dots + 1) % 4
			return m, spinnerTick()
		}
		return m, nil
	}

	mm, cmd := m.model.Update(msg)
	return streamModel{model: mm.(model), spinning: m.spinning, dots: m.dots}, cmd
}

func (m streamModel) View() string {
	// Reuse existing view, but add a small status line under the title.
	base := m.model.View()
	if !m.spinning {
		return base
	}
	dots := strings.Repeat(".", m.dots)
	status := "searching" + dots
	// Inject status after first line (title).
	lines := strings.SplitN(base, "\n", 2)
	if len(lines) < 2 {
		return base + "\n" + status
	}
	return lines[0] + "\n" + status + "\n" + lines[1]
}

func normalizeDevices(in []core.Device) []core.Device {
	// Keep unique by address, require name+address.
	m := map[string]core.Device{}
	for _, d := range in {
		addr := strings.TrimSpace(d.Address)
		if addr == "" || strings.TrimSpace(d.Name) == "" {
			continue
		}
		m[addr] = d
	}
	out := make([]core.Device, 0, len(m))
	for _, d := range m {
		out = append(out, d)
	}
	// Stable order for better UX.
	sort.SliceStable(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out
}

// PickDeviceStream opens a picker UI immediately and keeps updating the list
// with devices delivered by the updates channel until the user selects/cancels.
func (p Picker) PickDeviceStream(ctx context.Context, title string, updates <-chan []core.Device) (core.Device, error) {
	m := newStreamModel(title)
	program := tea.NewProgram(m, tea.WithContext(ctx), tea.WithAltScreen())

	// Fan-in updates to Bubble Tea.
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ds, ok := <-updates:
				if !ok {
					return
				}
				program.Send(devicesUpdateMsg{devices: ds})
			}
		}
	}()

	res, err := program.Run()
	if err != nil {
		return core.Device{}, err
	}
	fm, ok := res.(streamModel)
	if !ok {
		return core.Device{}, core.ErrCanceled{}
	}
	if fm.canceled {
		return core.Device{}, core.ErrCanceled{}
	}
	return fm.selected, nil
}
