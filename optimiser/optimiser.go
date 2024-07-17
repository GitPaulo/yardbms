package optimiser

import (
	"yardbms/models"
)

func OptimizeQuery(parsedQuery models.ParsedQuery) models.ParsedQuery {
	// For now, just return the parsed query without changes
	return parsedQuery
}
