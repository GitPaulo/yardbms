package storage

import (
	"fmt"
	"os"
	"sync"
	"yardbms/db/storage/locking"
	"yardbms/db/storage/transactions"
)

type RAMStorage struct {
	tables      map[string][]map[string]interface{}
	log         *transactions.TransactionLog
	lockManager *locking.LockManager
	lock        sync.Mutex
}

func NewRAMStorage() *RAMStorage {
	return &RAMStorage{
		tables:      make(map[string][]map[string]interface{}),
		log:         transactions.NewTransactionLog(),
		lockManager: locking.NewLockManager(),
	}
}

func (rs *RAMStorage) CreateTable(tableName string) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	if _, exists := rs.tables[tableName]; exists {
		return fmt.Errorf("table %s already exists", tableName)
	}
	rs.tables[tableName] = []map[string]interface{}{}
	return nil
}

func (rs *RAMStorage) DropTable(tableName string) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	if _, exists := rs.tables[tableName]; exists {
		delete(rs.tables, tableName)
		return nil
	}
	return fmt.Errorf("table %s does not exist", tableName)
}

func (rs *RAMStorage) Insert(tableName string, data map[string]interface{}, transactionID string) error {
	rs.lockManager.LockTable(tableName)
	defer rs.lockManager.UnlockTable(tableName)

	rs.lock.Lock()
	defer rs.lock.Unlock()

	if _, exists := rs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	rs.tables[tableName] = append(rs.tables[tableName], data)
	if transactionID != "" {
		rs.log.RecordOperation(transactionID, transactions.Operation{
			Type:      "INSERT",
			TableName: tableName,
			Row:       data,
		})
	}
	return nil
}

func (rs *RAMStorage) Select(tableName string) ([]map[string]interface{}, error) {
	rs.lockManager.RLockTable(tableName)
	defer rs.lockManager.RUnlockTable(tableName)

	rs.lock.Lock()
	defer rs.lock.Unlock()

	if data, exists := rs.tables[tableName]; exists {
		return data, nil
	}
	return nil, fmt.Errorf("table %s does not exist", tableName)
}

func (rs *RAMStorage) Update(tableName string, setClauses map[string]interface{}, whereClause string, transactionID string) error {
	rs.lockManager.LockTable(tableName)
	defer rs.lockManager.UnlockTable(tableName)

	rs.lock.Lock()
	defer rs.lock.Unlock()

	if _, exists := rs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	for i, row := range rs.tables[tableName] {
		if MatchesWhereClause(row, whereClause) {
			for col, value := range setClauses {
				oldValue := row[col]
				rs.tables[tableName][i][col] = value
				if transactionID != "" {
					rs.log.RecordOperation(transactionID, transactions.Operation{
						Type:      "UPDATE",
						TableName: tableName,
						Row:       map[string]interface{}{col: oldValue},
					})
				}
			}
		}
	}
	return nil
}

func (rs *RAMStorage) Delete(tableName string, whereClause string, transactionID string) error {
	rs.lockManager.LockTable(tableName)
	defer rs.lockManager.UnlockTable(tableName)

	rs.lock.Lock()
	defer rs.lock.Unlock()

	if _, exists := rs.tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	var newTable []map[string]interface{}
	for _, row := range rs.tables[tableName] {
		if !MatchesWhereClause(row, whereClause) {
			newTable = append(newTable, row)
		} else if transactionID != "" {
			rs.log.RecordOperation(transactionID, transactions.Operation{
				Type:      "DELETE",
				TableName: tableName,
				Row:       row,
			})
		}
	}
	rs.tables[tableName] = newTable
	return nil
}

func (rs *RAMStorage) StartTransaction(id string) {
	rs.log.StartTransaction(id)
}

func (rs *RAMStorage) CommitTransaction(id string) {
	rs.log.CommitTransaction(id)
	err := rs.log.SaveToDisk("transaction.log")
	if err != nil {
		fmt.Println("Error with commit transaction saving to disk:", err)
		os.Exit(1)
	}
}

func (rs *RAMStorage) RollbackTransaction(id string) {
	rs.log.RollbackTransaction(id, rs)
	err := rs.log.SaveToDisk("transaction.log")
	if err != nil {
		fmt.Println("Error with rollback transaction and saving to disk:", err)
		os.Exit(1)
	}
}

func (rs *RAMStorage) RollbackInsert(tableName string, row map[string]interface{}) {
	if table, exists := rs.tables[tableName]; exists {
		for i, r := range table {
			if MapsEqual(r, row) {
				rs.tables[tableName] = append(rs.tables[tableName][:i], rs.tables[tableName][i+1:]...)
				break
			}
		}
	}
}
