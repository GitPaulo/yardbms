package storage

import (
	"fmt"
	"strings"
)

type RAMStorage struct {
	tables map[string][]map[string]interface{}
}

func (rs *RAMStorage) CreateTable(tableName string) error {
	if _, exists := rs.tables[tableName]; exists {
		return fmt.Errorf("table %s already exists", tableName)
	}
	rs.tables[tableName] = []map[string]interface{}{}
	return nil
}

func (rs *RAMStorage) Insert(tableName string, data map[string]interface{}) error {
	if _, exists := rs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	rs.tables[tableName] = append(rs.tables[tableName], data)
	return nil
}

func (rs *RAMStorage) Select(tableName string) ([]map[string]interface{}, error) {
	if data, exists := rs.tables[tableName]; exists {
		return data, nil
	}
	return nil, fmt.Errorf("table %s does not exist", tableName)
}

func (rs *RAMStorage) Update(tableName string, setClauses map[string]interface{}, whereClause string) error {
	if _, exists := rs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	for i, row := range rs.tables[tableName] {
		if matchesWhereClauseRAM(row, whereClause) {
			for col, value := range setClauses {
				rs.tables[tableName][i][col] = value
			}
		}
	}
	return nil
}

func (rs *RAMStorage) Delete(tableName string, whereClause string) error {
	if _, exists := rs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	var newTable []map[string]interface{}
	for _, row := range rs.tables[tableName] {
		if !matchesWhereClauseRAM(row, whereClause) {
			newTable = append(newTable, row)
		}
	}
	rs.tables[tableName] = newTable
	return nil
}

func matchesWhereClauseRAM(row map[string]interface{}, whereClause string) bool {
	if whereClause == "" {
		return true
	}
	// TODO: Finish WHERE clause
	// Basic for now
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
