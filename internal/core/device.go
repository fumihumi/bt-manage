package core

import "time"

type Device struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Type    string `json:"type"`
	RSSI    *int   `json:"rssi,omitempty"`

	Connected       bool       `json:"connected"`
	LastConnectedAt *time.Time `json:"lastConnectedAt,omitempty"`
}
