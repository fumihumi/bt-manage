package blueutil

import (
	"testing"
)

func TestParseDeviceListJSON(t *testing.T) {
	in := []byte(`[
  {"address":"aa-bb-cc-dd-ee-ff","recentAccessDate":"2026-01-03T00:59:42+09:00","name":"MX Master","connected":true,"paired":true,"RSSI":-12},
  {"address":"11-22-33-44-55-66","recentAccessDate":"2026-01-02T00:00:00Z","name":"Keychron","connected":false,"paired":true}
]`)

	got, err := parseDeviceListJSON(in)
	if err != nil {
		t.Fatalf("parseDeviceListJSON: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2", len(got))
	}
	if got[0].Address != "aa:bb:cc:dd:ee:ff" {
		t.Fatalf("addr=%q", got[0].Address)
	}
	if got[0].RSSI == nil || *got[0].RSSI != -12 {
		t.Fatalf("rssi=%v", got[0].RSSI)
	}
	if got[1].RSSI != nil {
		t.Fatalf("rssi=%v", got[1].RSSI)
	}
}

func TestParseDeviceListJSON_RawRSSIOnly(t *testing.T) {
	in := []byte(`[
  {"address":"aa-bb-cc-dd-ee-ff","recentAccessDate":"2026-01-03T00:59:42+09:00","name":"MX Master","connected":true,"paired":true,"rawRSSI":-77}
]`)

	got, err := parseDeviceListJSON(in)
	if err != nil {
		t.Fatalf("parseDeviceListJSON: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len=%d, want 1", len(got))
	}
	if got[0].RSSI == nil || *got[0].RSSI != -77 {
		t.Fatalf("rssi=%v", got[0].RSSI)
	}
}

func TestParseDeviceListJSON_RSSIOnly(t *testing.T) {
	in := []byte(`[
  {"address":"aa-bb-cc-dd-ee-ff","recentAccessDate":"2026-01-03T00:59:42+09:00","name":"MX Master","connected":true,"paired":true,"RSSI":-55}
]`)

	got, err := parseDeviceListJSON(in)
	if err != nil {
		t.Fatalf("parseDeviceListJSON: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len=%d, want 1", len(got))
	}
	if got[0].RSSI == nil || *got[0].RSSI != -55 {
		t.Fatalf("rssi=%v", got[0].RSSI)
	}
}
