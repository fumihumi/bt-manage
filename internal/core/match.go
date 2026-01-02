package core

import "strings"

func findByName(devices []Device, query string, exact bool) []Device {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil
	}

	matches := make([]Device, 0)
	for _, d := range devices {
		if exact {
			if d.Name == q {
				matches = append(matches, d)
			}
			continue
		}

		if strings.HasPrefix(d.Name, q) {
			matches = append(matches, d)
		}
	}
	return matches
}
