package storage

type Storage interface {
	CreateTable(tableName string) error

	Insert(tableName string, data map[string]interface{}, transactionID string) error
	Select(tableName string) ([]map[string]interface{}, error)
	Update(tableName string, setClauses map[string]interface{}, whereClause string, transactionID string) error
	Delete(tableName string, whereClause string, transactionID string) error

	StartTransaction(id string)
	CommitTransaction(id string)
	RollbackTransaction(id string)
	rollbackInsert(tableName string, row map[string]interface{})
}
