package execution

import (
	"fmt"
	"yardbms/models"
	"yardbms/storage"

	"github.com/xwb1989/sqlparser"
)

func ExecuteQuery(optimizedQuery models.ParsedQuery, storage storage.Storage) string {
	switch stmt := optimizedQuery.Stmt.(type) {
	case *sqlparser.Select:
		return handleSelect(stmt, storage)
	case *sqlparser.DDL:
		return handleDDL(stmt, storage)
	case *sqlparser.Insert:
		return handleInsert(stmt, storage)
	case *sqlparser.Update:
		return handleUpdate(stmt, storage)
	case *sqlparser.Delete:
		return handleDelete(stmt, storage)
	default:
		return "Unsupported query type"
	}
}

func handleSelect(stmt *sqlparser.Select, storage storage.Storage) string {
	tableName := stmt.From[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
	rows, err := storage.Select(tableName)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return fmt.Sprintf("Rows: %v", rows)
}

func handleDDL(stmt *sqlparser.DDL, storage storage.Storage) string {
	if stmt.Action == "create" {
		err := storage.CreateTable(stmt.NewName.Name.String())
		if err != nil {
			return fmt.Sprintf("Error: %s", err)
		}
		return fmt.Sprintf("Table %s created", stmt.NewName.Name.String())
	}
	return "Unsupported DDL action"
}

func handleInsert(stmt *sqlparser.Insert, storage storage.Storage) string {
	tableName := stmt.Table.Name.String()
	row := make(map[string]interface{})
	for i, col := range stmt.Columns {
		row[col.String()] = sqlparser.String(stmt.Rows.(sqlparser.Values)[0][i])
	}
	err := storage.Insert(tableName, row)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return fmt.Sprintf("Row inserted into %s", tableName)
}

func handleUpdate(stmt *sqlparser.Update, storage storage.Storage) string {
	tableName := stmt.TableExprs[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
	setClauses := make(map[string]interface{})
	for _, expr := range stmt.Exprs {
		setClauses[expr.Name.Name.String()] = sqlparser.String(expr.Expr)
	}
	whereClause := sqlparser.String(stmt.Where)
	err := storage.Update(tableName, setClauses, whereClause)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return fmt.Sprintf("Table %s updated", tableName)
}

func handleDelete(stmt *sqlparser.Delete, storage storage.Storage) string {
	tableName := stmt.TableExprs[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
	whereClause := sqlparser.String(stmt.Where)
	err := storage.Delete(tableName, whereClause)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return fmt.Sprintf("Rows deleted from %s", tableName)
}
