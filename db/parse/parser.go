package parse

import (
	"yardbms/models"

	"github.com/xwb1989/sqlparser"
)

func ParseQuery(query string) (models.ParsedQuery, error) {
	stmt, err := sqlparser.Parse(query) // Thanks sqlparser! TODO: Write own parser (copium)
	if err != nil {
		return models.ParsedQuery{}, err
	}

	return models.ParsedQuery{Stmt: stmt}, nil
}
