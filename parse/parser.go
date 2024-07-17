package parser

import (
	"yardms/models"

	"github.com/xwb1989/sqlparser"
)

// ParseQuery uses sqlparser to parse the input query
func ParseQuery(query string) (models.ParsedQuery, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return models.ParsedQuery{}, err
	}

	return models.ParsedQuery{Stmt: stmt}, nil
}
