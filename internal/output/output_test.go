package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/fumihumi/bt-manage/internal/core"
)

func TestWriteTSV_WithHeader(t *testing.T) {
	rssi := -42
	devices := []core.Device{{Name: "MX", Address: "AA", Type: "Keyboard", RSSI: &rssi}}

	var buf bytes.Buffer
	if err := WriteTSV(&buf, devices, true); err != nil {
		t.Fatalf("WriteTSV: %v", err)
	}

	out := buf.String()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("lines=%d, want 2; out=%q", len(lines), out)
	}

	header := lines[0]
	for _, col := range []string{"Name", "Address", "Type", "RSSI"} {
		if !strings.Contains(header, col) {
			t.Fatalf("header missing %q: %q", col, header)
		}
	}
	if !strings.Contains(lines[1], "MX") || !strings.Contains(lines[1], "AA") || !strings.Contains(lines[1], "Keyboard") || !strings.Contains(lines[1], "-42") {
		t.Fatalf("row=%q", lines[1])
	}
}

func TestWriteTSV_NoHeader(t *testing.T) {
	devices := []core.Device{{Name: "MX", Address: "AA", Type: "", RSSI: nil}}

	var buf bytes.Buffer
	if err := WriteTSV(&buf, devices, false); err != nil {
		t.Fatalf("WriteTSV: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	if strings.Contains(out, "Name") && strings.Contains(out, "Address") && strings.Contains(out, "Type") && strings.Contains(out, "RSSI") {
		t.Fatalf("unexpected header-like output: %q", out)
	}
}

func TestWriteJSON(t *testing.T) {
	devices := []core.Device{{Name: "MX", Address: "AA", Type: "Keyboard"}}

	var buf bytes.Buffer
	if err := WriteJSON(&buf, devices); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}

	var got []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len=%d, want 1", len(got))
	}
	if got[0]["name"] != "MX" {
		t.Fatalf("name=%v", got[0]["name"])
	}
}
