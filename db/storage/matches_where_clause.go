package storage

import (
	"fmt"
	"strings"
)

func MatchesWhereClause(row map[string]interface{}, whereClause string) bool {
	if whereClause == "" {
		return true
	}
	conditions := strings.Split(whereClause, "AND")
	for _, condition := range conditions {
		condition = strings.TrimSpace(condition)
		parts := strings.Split(condition, "=")
		if len(parts) != 2 {
			continue
		}
		col := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if fmt.Sprintf("%v", row[col]) != val {
			return false
		}
	}
	return true
}
