package engine

import (
	"fmt"
	"strconv"
	"strings"
	"yardbms/db/models"

	"github.com/xwb1989/sqlparser"
)

func ExecuteQuery(optimizedQuery models.ParsedQuery, storage models.Storage, transactionID string) string {
	switch stmt := optimizedQuery.Stmt.(type) {
	case *sqlparser.Select:
		return handleSelect(stmt, storage)
	case *sqlparser.DDL:
		return handleDDL(stmt, storage)
	case *sqlparser.Insert:
		return handleInsert(stmt, storage, transactionID)
	case *sqlparser.Update:
		return handleUpdate(stmt, storage, transactionID)
	case *sqlparser.Delete:
		return handleDelete(stmt, storage, transactionID)
	default:
		return "Unsupported query type"
	}
}

func handleSelect(stmt *sqlparser.Select, storage models.Storage) string {
	tableName := stmt.From[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
	rows, err := storage.Select(tableName)
	if err != nil {
		return formatError(err)
	}

	selectExprs := stmt.SelectExprs
	if len(selectExprs) == 1 {
		if _, ok := selectExprs[0].(*sqlparser.StarExpr); ok {
			return formatSelectResult(rows)
		}
	}

	return handleAggregateFunctions(selectExprs, rows)
}

func handleAggregateFunctions(selectExprs sqlparser.SelectExprs, rows []map[string]interface{}) string {
	results := make(map[string]interface{})
	for _, expr := range selectExprs {
		switch e := expr.(type) {
		case *sqlparser.AliasedExpr:
			switch agg := e.Expr.(type) {
			case *sqlparser.FuncExpr:
				funcName := strings.ToLower(agg.Name.String())
				var column string
				if len(agg.Exprs) > 0 {
					switch colExpr := agg.Exprs[0].(type) {
					case *sqlparser.AliasedExpr:
						column = colExpr.Expr.(*sqlparser.ColName).Name.String()
					case *sqlparser.StarExpr:
						column = "*"
					}
				}
				switch funcName {
				case "count":
					results[funcName] = len(rows)
				case "sum":
					results[funcName] = calculateSum(rows, column)
				case "avg":
					results[funcName] = calculateAverage(rows, column)
				case "min":
					results[funcName] = findMin(rows, column)
				case "max":
					results[funcName] = findMax(rows, column)
				}
			}
		}
	}

	return formatAggregateResult(results)
}

func calculateSum(rows []map[string]interface{}, column string) float64 {
	sum := 0.0
	for _, row := range rows {
		value, _ := strconv.ParseFloat(fmt.Sprintf("%v", row[column]), 64)
		sum += value
	}
	return sum
}

func calculateAverage(rows []map[string]interface{}, column string) float64 {
	sum := calculateSum(rows, column)
	return sum / float64(len(rows))
}

func findMin(rows []map[string]interface{}, column string) float64 {
	if len(rows) == 0 {
		return 0
	}
	min, _ := strconv.ParseFloat(fmt.Sprintf("%v", rows[0][column]), 64)
	for _, row := range rows {
		value, _ := strconv.ParseFloat(fmt.Sprintf("%v", row[column]), 64)
		if value < min {
			min = value
		}
	}
	return min
}

func findMax(rows []map[string]interface{}, column string) float64 {
	if len(rows) == 0 {
		return 0
	}
	max, _ := strconv.ParseFloat(fmt.Sprintf("%v", rows[0][column]), 64)
	for _, row := range rows {
		value, _ := strconv.ParseFloat(fmt.Sprintf("%v", row[column]), 64)
		if value > max {
			max = value
		}
	}
	return max
}

func handleDDL(stmt *sqlparser.DDL, storage models.Storage) string {
	switch stmt.Action {
	case "create":
		err := storage.CreateTable(stmt.NewName.Name.String())
		if err != nil {
			return formatError(err)
		}
		return formatSuccess(fmt.Sprintf("Table %s created", stmt.NewName.Name.String()))
	case "drop":
		err := storage.DropTable(stmt.NewName.Name.String())
		if err != nil {
			return formatError(err)
		}
		return formatSuccess(fmt.Sprintf("Table %s dropped", stmt.NewName.Name.String()))
	default:
		return "Unsupported DDL action"
	}
}

func handleInsert(stmt *sqlparser.Insert, storage models.Storage, transactionID string) string {
	tableName := stmt.Table.Name.String()
	row := make(map[string]interface{})
	for i, col := range stmt.Columns {
		row[col.String()] = sqlparser.String(stmt.Rows.(sqlparser.Values)[0][i])
	}
	err := storage.Insert(tableName, row, transactionID)
	if err != nil {
		return formatError(err)
	}
	return formatSuccess(fmt.Sprintf("Row inserted into %s", tableName))
}

func handleUpdate(stmt *sqlparser.Update, storage models.Storage, transactionID string) string {
	tableName := stmt.TableExprs[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
	setClauses := make(map[string]interface{})
	for _, expr := range stmt.Exprs {
		setClauses[expr.Name.Name.String()] = sqlparser.String(expr.Expr)
	}
	whereClause := sqlparser.String(stmt.Where)
	err := storage.Update(tableName, setClauses, whereClause, transactionID)
	if err != nil {
		return formatError(err)
	}
	return formatSuccess(fmt.Sprintf("Table %s updated", tableName))
}

func handleDelete(stmt *sqlparser.Delete, storage models.Storage, transactionID string) string {
	tableName := stmt.TableExprs[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
	whereClause := sqlparser.String(stmt.Where)
	err := storage.Delete(tableName, whereClause, transactionID)
	if err != nil {
		return formatError(err)
	}
	return formatSuccess(fmt.Sprintf("Rows deleted from %s", tableName))
}

func formatSelectResult(rows []map[string]interface{}) string {
	if len(rows) == 0 {
		return "No rows found"
	}
	result := "Rows:\n"
	for _, row := range rows {
		for col, val := range row {
			result += fmt.Sprintf("├─ %s: %v\n", col, val)
		}
		result += "└─\n"
	}
	return result
}

func formatAggregateResult(results map[string]interface{}) string {
	result := "Aggregate Results:\n"
	for funcName, value := range results {
		result += fmt.Sprintf("├─ %s: %v\n", strings.Title(funcName), value)
	}
	result += "└─\n"
	return result
}

func formatError(err error) string {
	return fmt.Sprintf("Error: %s\n", err)
}

func formatSuccess(message string) string {
	return fmt.Sprintf("Success: %s\n", message)
}
