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

func TestMultiModelToggleAndConfirm(t *testing.T) {
	devices := []core.Device{
		{Name: "A", Address: "AA", Connected: true},
		{Name: "B", Address: "BB", Connected: false},
		{Name: "C", Address: "CC", Connected: false},
	}

	m := newMultiModel("Connect", devices)
	if len(m.filtered) != 2 {
		t.Fatalf("expected connected device to be filtered out; filtered=%d", len(m.filtered))
	}

	// Toggle first visible item (B).
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = mm.(multiModel)
	if !m.selectedMap["BB"] {
		t.Fatalf("expected BB to be selected")
	}

	// Move down and toggle second (C).
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = mm.(multiModel)
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = mm.(multiModel)
	if !m.selectedMap["CC"] {
		t.Fatalf("expected CC to be selected")
	}

	// Confirm.
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mm.(multiModel)
	if !m.chosen || m.canceled {
		t.Fatalf("expected chosen=true canceled=false")
	}
	if len(m.selected) != 2 {
		t.Fatalf("expected 2 selected devices, got %d", len(m.selected))
	}
}

func TestMultiModelCancel(t *testing.T) {
	m := newMultiModel("Connect", []core.Device{{Name: "A", Address: "AA"}})
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = mm.(multiModel)
	if !m.canceled {
		t.Fatalf("expected canceled=true")
	}
}
