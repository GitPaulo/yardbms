package db

import (
	"yardbms/db/engine"
	"yardbms/db/optimiser"
	"yardbms/db/parse"
	"yardbms/storage"
)

type Database struct {
	storage storage.Storage
}

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

func (db *Database) ExecuteQuery(query string, transactionID string) (string, error) {
	// Parse the query
	parsedQuery, err := parse.ParseQuery(query)
	if err != nil {
		return "", err
	}

	// Optimize the query
	optimizedQuery := optimiser.OptimizeQuery(parsedQuery)

	// Execute the query
	result := engine.ExecuteQuery(optimizedQuery, db.storage, transactionID)

	return result, nil
}

func (db *Database) StartTransaction(transactionID string) {
	db.storage.StartTransaction(transactionID)
}

func (db *Database) CommitTransaction(transactionID string) {
	db.storage.CommitTransaction(transactionID)
}

func (db *Database) RollbackTransaction(transactionID string) {
	db.storage.RollbackTransaction(transactionID)
}
