package db

import (
	"fmt"
	"time"
)

type TransactionManager struct {
	db *Database
}

func NewTransactionManager(db *Database) *TransactionManager {
	return &TransactionManager{db: db}
}

func (tm *TransactionManager) StartTransaction() string {
	transactionID := fmt.Sprintf("txn_%d", time.Now().UnixNano())
	tm.db.StartTransaction(transactionID)
	return transactionID
}

func (tm *TransactionManager) CommitTransaction(transactionID string) error {
	if transactionID == "" {
		return fmt.Errorf("no active transaction to commit")
	}
	tm.db.CommitTransaction(transactionID)
	return nil
}

func (tm *TransactionManager) RollbackTransaction(transactionID string) error {
	if transactionID == "" {
		return fmt.Errorf("no active transaction to rollback")
	}
	tm.db.RollbackTransaction(transactionID)
	return nil
}
