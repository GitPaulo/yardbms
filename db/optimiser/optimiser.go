package optimiser

import (
	"yardbms/db/models"

	"github.com/xwb1989/sqlparser"
)

func OptimizeQuery(parsedQuery models.ParsedQuery) models.ParsedQuery {
	switch stmt := parsedQuery.Stmt.(type) {
	case *sqlparser.Select:
		optimizeSelect(stmt)
	}
	return parsedQuery
}

func optimizeSelect(sel *sqlparser.Select) {
	predicatePushdown(sel)
	eliminateRedundantProjections(sel)
}

func predicatePushdown(sel *sqlparser.Select) {
	pushdownWhere(sel.From, sel.Where)
}

func pushdownWhere(tableExprs sqlparser.TableExprs, where *sqlparser.Where) {
	// Traverse the table expressions to find join conditions and subqueries
	for _, tableExpr := range tableExprs {
		switch tableExpr := tableExpr.(type) {
		case *sqlparser.AliasedTableExpr:
			if subquery, ok := tableExpr.Expr.(*sqlparser.Subquery); ok {
				// If it's a subquery, push the WHERE clause into the subquery
				if sel, ok := subquery.Select.(*sqlparser.Select); ok {
					sel.Where = combineWhere(sel.Where, where)
				}
			}
		case *sqlparser.JoinTableExpr:
			// Handle JOINs by pushing down the WHERE clause into the joined tables
			pushdownWhereForTableExpr(tableExpr.LeftExpr, where)
			pushdownWhereForTableExpr(tableExpr.RightExpr, where)
		}
	}
}

func pushdownWhereForTableExpr(tableExpr sqlparser.TableExpr, where *sqlparser.Where) {
	switch tableExpr := tableExpr.(type) {
	case *sqlparser.AliasedTableExpr:
		if subquery, ok := tableExpr.Expr.(*sqlparser.Subquery); ok {
			// If it's a subquery, push the WHERE clause into the subquery
			if sel, ok := subquery.Select.(*sqlparser.Select); ok {
				sel.Where = combineWhere(sel.Where, where)
			}
		}
	case *sqlparser.JoinTableExpr:
		// Handle JOINs by pushing down the WHERE clause into the joined tables
		pushdownWhereForTableExpr(tableExpr.LeftExpr, where)
		pushdownWhereForTableExpr(tableExpr.RightExpr, where)
	}
}

func combineWhere(existingWhere, newWhere *sqlparser.Where) *sqlparser.Where {
	if existingWhere == nil {
		return newWhere
	}
	if newWhere == nil {
		return existingWhere
	}

	// Combine existing WHERE and new WHERE using AND
	return &sqlparser.Where{
		Type: sqlparser.WhereStr,
		Expr: &sqlparser.AndExpr{
			Left:  existingWhere.Expr,
			Right: newWhere.Expr,
		},
	}
}

func eliminateRedundantProjections(sel *sqlparser.Select) {
	neededColumns := map[string]bool{
		"id":   true,
		"name": true,
	}
	newSelectExprs := sqlparser.SelectExprs{}
	for _, selectExpr := range sel.SelectExprs {
		switch expr := selectExpr.(type) {
		case *sqlparser.AliasedExpr:
			if colName, ok := expr.Expr.(*sqlparser.ColName); ok {
				column := colName.Name.String()
				if neededColumns[column] {
					newSelectExprs = append(newSelectExprs, selectExpr)
				}
			}
		default:
			newSelectExprs = append(newSelectExprs, selectExpr)
		}
	}
	sel.SelectExprs = newSelectExprs
}
