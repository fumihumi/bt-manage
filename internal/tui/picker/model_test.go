package picker

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fumihumi/bt-manage/internal/core"
)

func TestModel_Cancel(t *testing.T) {
	m := newModel("Pick", []core.Device{{Name: "A"}})
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m2 := mm.(model)
	if !m2.canceled {
		t.Fatalf("expected canceled")
	}
}

func TestModel_Filter_DownUp(t *testing.T) {
	m := newModel("Pick", []core.Device{{Name: "Alpha"}, {Name: "Beta"}, {Name: "Gamma"}})

	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m2 := mm.(model)
	if m2.index != 1 {
		t.Fatalf("index=%d", m2.index)
	}

	mm, _ = m2.Update(tea.KeyMsg{Type: tea.KeyUp})
	m3 := mm.(model)
	if m3.index != 0 {
		t.Fatalf("index=%d", m3.index)
	}

	// filter by typing "ga"
	m3.input.SetValue("ga")
	m3.applyFilter()
	if len(m3.filtered) != 1 || m3.filtered[0].Name != "Gamma" {
		t.Fatalf("filtered=%v", m3.filtered)
	}
}
