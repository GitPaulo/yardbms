package main

import (
	"fmt"
	"yardms/execution"
	"yardms/input"
	"yardms/optimiser"
	"yardms/parse"
)

func main() {
	// Sample query input
	query := input.GetQuery()

	// Parse the query
	parsedQuery, err := parser.ParseQuery(query)
	if err != nil {
		fmt.Println("Error parsing query:", err)
		return
	}

	// Optimize the query
	optimizedQuery := optimizer.OptimizeQuery(parsedQuery)

	// Execute the query
	result := execution.ExecuteQuery(optimizedQuery)

	// Print the result
	fmt.Println("Query Result:", result)
}
