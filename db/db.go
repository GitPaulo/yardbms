package db

import (
	"yardbms/execution"
	"yardbms/optimiser"
	"yardbms/parse"
	"yardbms/storage"
)

type Database struct {
	storage storage.Storage
}

func New(storageType string) *Database {
	var st storage.Storage
	if storageType == "file" {
		st = storage.NewFileStorage()
	} else {
		st = storage.NewRAMStorage() // Default to RAM storage
	}

	return &Database{
		storage: st,
	}
}

func (db *Database) ExecuteQuery(query string) (string, error) {
	// Parse the query
	parsedQuery, err := parse.ParseQuery(query)
	if err != nil {
		return "", err
	}

	// Optimize the query
	optimizedQuery := optimiser.OptimizeQuery(parsedQuery)

	// Execute the query
	result := execution.ExecuteQuery(optimizedQuery, db.storage)

	return result, nil
}
