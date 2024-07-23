package transactions

import (
	"encoding/json"
	"os"
	"yardbms/db/models"
)

type TransactionLog struct {
	Transactions map[string]TransactionRecord
}

type TransactionRecord struct {
	TransactionID string
	Operations    []Operation
}

type Operation struct {
	Type      string
	TableName string
	Row       map[string]interface{}
}

func NewTransactionLog() *TransactionLog {
	return &TransactionLog{
		Transactions: make(map[string]TransactionRecord),
	}
}

func (log *TransactionLog) StartTransaction(id string) {
	log.Transactions[id] = TransactionRecord{
		TransactionID: id,
		Operations:    []Operation{},
	}
}

func (log *TransactionLog) RecordOperation(id string, op Operation) {
	if record, exists := log.Transactions[id]; exists {
		record.Operations = append(record.Operations, op)
		log.Transactions[id] = record
	}
}

func (log *TransactionLog) CommitTransaction(id string) {
	delete(log.Transactions, id)
}

func (log *TransactionLog) RollbackTransaction(id string, storage models.Storage) {
	if record, exists := log.Transactions[id]; exists {
		for _, op := range record.Operations {
			switch op.Type {
			case "INSERT":
				storage.RollbackInsert(op.TableName, op.Row)
			}
		}
		delete(log.Transactions, id)
	}
}

func (log *TransactionLog) SaveToDisk(filePath string) error {
	data, err := json.Marshal(log)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func (log *TransactionLog) LoadFromDisk(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, log)
}
