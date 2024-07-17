package parse

import (
	"yardbms/models"

	"github.com/xwb1989/sqlparser"
)

func ParseQuery(query string) (models.ParsedQuery, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return models.ParsedQuery{}, err
	}

	return models.ParsedQuery{Stmt: stmt}, nil
}
