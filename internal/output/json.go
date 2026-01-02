package output

import (
	"encoding/json"
	"io"

	"github.com/fumihumi/bt-manage/internal/core"
)

func WriteJSON(w io.Writer, devices []core.Device) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(devices)
}
