package blueutil

import (
	"encoding/json"
	"time"

	"github.com/fumihumi/bt-manage/internal/core"
)

type deviceJSON struct {
	Address          string     `json:"address"`
	Name             string     `json:"name"`
	Connected        bool       `json:"connected"`
	RecentAccessDate *time.Time `json:"recentAccessDate"`
	RSSI             *int       `json:"RSSI"`
	RawRSSI          *int       `json:"rawRSSI"`
	Paired           bool       `json:"paired"`
}

func parseDeviceListJSON(b []byte) ([]core.Device, error) {
	var in []deviceJSON
	if err := json.Unmarshal(b, &in); err != nil {
		return nil, err
	}

	out := make([]core.Device, 0, len(in))
	for _, d := range in {
		dev := core.Device{
			Name:            d.Name,
			Address:         normalizeAddress(d.Address),
			Connected:       d.Connected,
			LastConnectedAt: d.RecentAccessDate,
			RSSI:            firstNonNilInt(d.RSSI, d.RawRSSI),
			Type:            "",
		}
		out = append(out, dev)
	}
	return out, nil
}

func normalizeAddress(addr string) string {
	// bt-manage では表示上は ':' 区切りに寄せる
	// (blueutil は '-' で返す)
	b := []byte(addr)
	for i := range b {
		if b[i] == '-' {
			b[i] = ':'
		}
	}
	return string(b)
}

func firstNonNilInt(a, b *int) *int {
	if a != nil {
		return a
	}
	return b
}
