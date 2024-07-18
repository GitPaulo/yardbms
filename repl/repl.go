package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"yardbms/db"
)

func Start(storageType string, filePath string) {
	reader := bufio.NewReader(os.Stdin)
	yardbms := db.New(storageType, filePath)
	transactionManager := db.NewTransactionManager(yardbms)

	fmt.Println("Welcome to the yardbms REPL")
	fmt.Println("Type 'exit' to quit")

	var transactionID string
	for {
		fmt.Print("yardbms> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		if strings.HasPrefix(input, "BEGIN TRANSACTION") {
			transactionID = transactionManager.StartTransaction()
			fmt.Println("Transaction started:", transactionID)
			continue
		}

		if strings.HasPrefix(input, "COMMIT") {
			err := transactionManager.CommitTransaction(transactionID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Transaction committed:", transactionID)
				transactionID = ""
			}
			continue
		}

		if strings.HasPrefix(input, "ROLLBACK") {
			err := transactionManager.RollbackTransaction(transactionID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Transaction rolled back:", transactionID)
				transactionID = ""
			}
			continue
		}

		result, err := yardbms.ExecuteQuery(input, transactionID)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Query Result:", result)
		}
	}
}
