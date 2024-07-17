package input

import (
	"bufio"
	"fmt"
)
import "os"

// GetQuery reads a SQL query from standard input
func GetQuery() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter SQL query: ")
	query, _ := reader.ReadString('\n')
	return query
}
