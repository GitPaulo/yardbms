package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"yardbms/db/storage/locking"
	"yardbms/db/storage/transactions"
)

type FileStorage struct {
	filePath    string
	tables      map[string][]map[string]interface{}
	log         *transactions.TransactionLog
	lockManager *locking.LockManager
}

func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		filePath:    filePath,
		tables:      make(map[string][]map[string]interface{}),
		log:         transactions.NewTransactionLog(),
		lockManager: locking.NewLockManager(),
	}
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

func (fs *FileStorage) DropTable(tableName string) error {
	if err := fs.load(); err != nil {
		return err
	}
	if _, exists := fs.tables[tableName]; exists {
		delete(fs.tables, tableName)
		return fs.save()
	}
	return fmt.Errorf("table %s does not exist", tableName)
}

func (fs *FileStorage) Insert(tableName string, data map[string]interface{}, transactionID string) error {
	if err := fs.load(); err != nil {
		return err
	}
	if _, exists := fs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	fs.tables[tableName] = append(fs.tables[tableName], data)
	if transactionID != "" {
		fs.log.RecordOperation(transactionID, transactions.Operation{
			Type:      "INSERT",
			TableName: tableName,
			Row:       data,
		})
	}
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

func (fs *FileStorage) Update(tableName string, setClauses map[string]interface{}, whereClause string, transactionID string) error {
	if err := fs.load(); err != nil {
		return err
	}
	if _, exists := fs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	for i, row := range fs.tables[tableName] {
		if MatchesWhereClause(row, whereClause) {
			for col, value := range setClauses {
				oldValue := row[col]
				fs.tables[tableName][i][col] = value
				if transactionID != "" {
					fs.log.RecordOperation(transactionID, transactions.Operation{
						Type:      "UPDATE",
						TableName: tableName,
						Row:       map[string]interface{}{col: oldValue},
					})
				}
			}
		}
	}
	return fs.save()
}

func (fs *FileStorage) Delete(tableName string, whereClause string, transactionID string) error {
	if err := fs.load(); err != nil {
		return err
	}
	if _, exists := fs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	var newTable []map[string]interface{}
	for _, row := range fs.tables[tableName] {
		if !MatchesWhereClause(row, whereClause) {
			newTable = append(newTable, row)
		} else if transactionID != "" {
			fs.log.RecordOperation(transactionID, transactions.Operation{
				Type:      "DELETE",
				TableName: tableName,
				Row:       row,
			})
		}
	}
	fs.tables[tableName] = newTable
	return fs.save()
}

func (fs *FileStorage) StartTransaction(id string) {
	fs.log.StartTransaction(id)
}

func (fs *FileStorage) CommitTransaction(id string) {
	fs.log.CommitTransaction(id)
	fs.log.SaveToDisk("transaction.log")
}

func (fs *FileStorage) RollbackTransaction(id string) {
	fs.log.RollbackTransaction(id, fs)
	fs.log.SaveToDisk("transaction.log")
}

func (fs *FileStorage) RollbackInsert(tableName string, row map[string]interface{}) {
	if table, exists := fs.tables[tableName]; exists {
		for i, r := range table {
			if MapsEqual(r, row) {
				fs.tables[tableName] = append(fs.tables[tableName][:i], fs.tables[tableName][i+1:]...)
				break
			}
		}
	}
}
