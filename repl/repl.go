package repl

import (
	"bufio"
	"fmt"
	"os"
	"yardms/db"
)

func Start(storageType string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the RDBMS REPL")
	fmt.Println("Type 'exit' to quit")

	db := db.New(storageType)

	for {
		fmt.Print("yardms> ")
		input, _ := reader.ReadString('\n')
		if input == "exit\n" {
			break
		}

		result, err := db.ExecuteQuery(input)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Query Result:", result)
		}
	}
}
