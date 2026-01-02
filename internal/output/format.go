package output

import (
	"fmt"
	"strings"
)

type Format int

const (
	FormatTSV Format = iota
	FormatJSON
)

func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "tsv":
		return FormatTSV, nil
	case "json":
		return FormatJSON, nil
	default:
		return 0, fmt.Errorf("unknown format: %s", s)
	}
}
