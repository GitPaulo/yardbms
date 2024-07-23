package parse

import (
	"yardbms/db/models"

	"github.com/xwb1989/sqlparser"
)

func ParseQuery(query string) (models.ParsedQuery, error) {
	// I bet you where expecting a custom parser....
	// Absolutely not: https://pkg.go.dev/github.com/xwb1989/sqlparser#section-documentation
	stmt, err := sqlparser.Parse(query)

	if err != nil {
		return models.ParsedQuery{}, err
	}

	return models.ParsedQuery{Stmt: stmt}, nil
}
