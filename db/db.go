package db

import (
	"yardbms/db/engine"
	"yardbms/db/optimiser"
	"yardbms/db/parse"
	"yardbms/storage"
)

// Database represents a database instance
type Database struct {
	storage storage.Storage
}

// New initializes a new database instance
func New(storageType string, filePath string) *Database {
	var st storage.Storage
	if storageType == "file" {
		st = storage.NewFileStorage(filePath)
	} else {
		st = storage.NewRAMStorage() // Default to RAM storage
	}

	return &Database{
		storage: st,
	}
}

// ExecuteQuery executes a given SQL query on the database
// Var args on transactionID is to allow for optional transaction ID
func (db *Database) ExecuteQuery(query string, transactionID ...string) (string, error) {
	// Parse the query
	parsedQuery, err := parse.ParseQuery(query)
	if err != nil {
		return "", err
	}

	// Optimize the query
	optimizedQuery := optimiser.OptimizeQuery(parsedQuery)

	// Determine the transaction ID to use
	var txnID string
	if len(transactionID) > 0 {
		txnID = transactionID[0]
	}

	// Execute the query
	result := engine.ExecuteQuery(optimizedQuery, db.storage, txnID)

	return result, nil
}

// StartTransaction starts a new transaction
func (db *Database) StartTransaction(transactionID string) {
	db.storage.StartTransaction(transactionID)
}

// CommitTransaction commits the current transaction
func (db *Database) CommitTransaction(transactionID string) {
	db.storage.CommitTransaction(transactionID)
}

// RollbackTransaction rolls back the current transaction
func (db *Database) RollbackTransaction(transactionID string) {
	db.storage.RollbackTransaction(transactionID)
}
