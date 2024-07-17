package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type FileStorage struct {
	filePath string
	tables   map[string][]map[string]interface{}
}

func (fs *FileStorage) load() error {
	file, err := os.OpenFile(fs.filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		return json.Unmarshal(data, &fs.tables)
	}

	fs.tables = make(map[string][]map[string]interface{})
	return nil
}

func (fs *FileStorage) save() error {
	file, err := os.OpenFile(fs.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(fs.tables)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func (fs *FileStorage) CreateTable(tableName string) error {
	if err := fs.load(); err != nil {
		return err
	}
	if _, exists := fs.tables[tableName]; exists {
		return fmt.Errorf("table %s already exists", tableName)
	}
	fs.tables[tableName] = []map[string]interface{}{}
	return fs.save()
}

func (fs *FileStorage) Insert(tableName string, data map[string]interface{}) error {
	if err := fs.load(); err != nil {
		return err
	}
	if _, exists := fs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	fs.tables[tableName] = append(fs.tables[tableName], data)
	return fs.save()
}

func (fs *FileStorage) Select(tableName string) ([]map[string]interface{}, error) {
	if err := fs.load(); err != nil {
		return nil, err
	}
	if data, exists := fs.tables[tableName]; exists {
		return data, nil
	}
	return nil, fmt.Errorf("table %s does not exist", tableName)
}

func (fs *FileStorage) Update(tableName string, setClauses map[string]interface{}, whereClause string) error {
	if err := fs.load(); err != nil {
		return err
	}
	if _, exists := fs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	for i, row := range fs.tables[tableName] {
		if matchesWhereClauseFile(row, whereClause) {
			for col, value := range setClauses {
				fs.tables[tableName][i][col] = value
			}
		}
	}
	return fs.save()
}

func (fs *FileStorage) Delete(tableName string, whereClause string) error {
	if err := fs.load(); err != nil {
		return err
	}
	if _, exists := fs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	var newTable []map[string]interface{}
	for _, row := range fs.tables[tableName] {
		if !matchesWhereClauseFile(row, whereClause) {
			newTable = append(newTable, row)
		}
	}
	fs.tables[tableName] = newTable
	return fs.save()
}

func matchesWhereClauseFile(row map[string]interface{}, whereClause string) bool {
	if whereClause == "" {
		return true
	}
	// TODO: Finish WHERE Clause
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
