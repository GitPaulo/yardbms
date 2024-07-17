package storage

type Storage interface {
	CreateTable(tableName string) error
	Insert(tableName string, data map[string]interface{}) error
	Select(tableName string) ([]map[string]interface{}, error)
	Update(tableName string, setClauses map[string]interface{}, whereClause string) error
	Delete(tableName string, whereClause string) error
}

func NewRAMStorage() Storage {
	return &RAMStorage{
		tables: make(map[string][]map[string]interface{}),
	}
}

func NewFileStorage() Storage {
	return &FileStorage{
		filePath: "data.json",
	}
}
