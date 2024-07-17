package models

import "github.com/xwb1989/sqlparser"

// ParsedQuery is a wrapper for the parsed SQL statement
type ParsedQuery struct {
	Stmt sqlparser.Statement
}

func (pq ParsedQuery) String() string {
	return sqlparser.String(pq.Stmt)
}
