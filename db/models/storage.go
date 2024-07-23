package models

type Storage interface {
	CreateTable(tableName string) error
	DropTable(tableName string) error

	Insert(tableName string, data map[string]interface{}, transactionID string) error
	Select(tableName string) ([]map[string]interface{}, error)
	Update(tableName string, setClauses map[string]interface{}, whereClause string, transactionID string) error
	Delete(tableName string, whereClause string, transactionID string) error

	StartTransaction(id string)
	CommitTransaction(id string)
	RollbackTransaction(id string)
	RollbackInsert(tableName string, row map[string]interface{})
}
