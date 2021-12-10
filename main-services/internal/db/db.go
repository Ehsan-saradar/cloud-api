package db

import (
	"strings"
)

// Helper function to join posibbly empty filters for a WHERE clause.
// Empty strings are discarded.
func Where(filters ...string) string {
	actualFilters := []string{}
	for _, filter := range filters {
		if filter != "" {
			actualFilters = append(actualFilters, filter)
		}
	}
	if len(actualFilters) == 0 {
		return ""
	}
	return "WHERE (" + strings.Join(actualFilters, ") AND (") + ")"
}
