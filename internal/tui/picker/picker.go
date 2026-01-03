package picker

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fumihumi/bt-manage/internal/core"
)

type Picker struct{}

func (p Picker) PickDevice(ctx context.Context, title string, devices []core.Device) (core.Device, error) {
	m := newModel(title, devices)

	program := tea.NewProgram(m, tea.WithContext(ctx), tea.WithAltScreen())
	res, err := program.Run()
	if err != nil {
		return core.Device{}, err
	}

	fm, ok := res.(model)
	if !ok {
		return core.Device{}, core.ErrCanceled{}
	}
	if fm.canceled {
		return core.Device{}, core.ErrCanceled{}
	}
	return fm.selected, nil
}

func (p Picker) PickDevices(ctx context.Context, title string, devices []core.Device) ([]core.Device, error) {
	m := newMultiModel(title, devices)

	program := tea.NewProgram(m, tea.WithContext(ctx), tea.WithAltScreen())
	res, err := program.Run()
	if err != nil {
		return nil, err
	}

	fm, ok := res.(multiModel)
	if !ok {
		return nil, core.ErrCanceled{}
	}
	if fm.canceled {
		return nil, core.ErrCanceled{}
	}
	return fm.selected, nil
}
