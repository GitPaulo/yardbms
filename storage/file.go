package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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
