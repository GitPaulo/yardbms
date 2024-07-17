package storage

import "fmt"

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
