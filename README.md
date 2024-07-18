# yardbms
Yes A/nother Relational Database Management System.

```sh
yardbms --storage=ram
```

In the REPL,

```sql
yardbms> CREATE TABLE users;
Table users created
yardbms> INSERT INTO users (id, name) VALUES (1, 'John Doe');
Row inserted into users
yardbms> SELECT * FROM users;
Rows: [map[id:1 name:'John Doe']]
```

As a package,

```go

import (
  "fmt"
  "yardbms/db"
)

func main() {
  // Initialize the database with file storage, for RAM use "ram" and ""
  yardbms := db.New("file", "file_path.json")

  // Create a table
  result, err := yardbms.ExecuteQuery("CREATE TABLE Users (id INT, name TEXT);")
  if err != nil {
    fmt.Println("Error:", err)
    return
  }
  fmt.Println("Query Result:", result)

  // Insert a row into the table
  result, err = yardbms.ExecuteQuery("INSERT INTO Users (id, name) VALUES (1, 'Alice');")
  if err != nil {
    fmt.Println("Error:", err)
    return
  }
  fmt.Println("Query Result:", result)

  // Select rows from the table
  result, err = yardbms.ExecuteQuery("SELECT * FROM Users;", "")
  if err != nil {
    fmt.Println("Error:", err)
    return
  }
  fmt.Println("Query Result:", result)

  // Start a transaction
  transactionManager := yardbms.NewTransactionManager(yardbms)
  transactionID := transactionManager.StartTransaction()
  fmt.Println("Transaction started:", transactionID)

  // Update a row within the transaction
  result, err = yardbms.ExecuteQuery("UPDATE Users SET name = 'Bob' WHERE id = 1;", transactionID)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }
  fmt.Println("Query Result:", result)

  // Commit the transaction
  err = transactionManager.CommitTransaction(transactionID)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }
  fmt.Println("Transaction committed:", transactionID)

  // Select rows from the table after commit
  result, err = yardbms.ExecuteQuery("SELECT * FROM Users;")
  if err != nil {
    fmt.Println("Error:", err)
    return
  }
  fmt.Println("Query Result:", result)
}
```

# Some plan IDK

![image](https://github.com/user-attachments/assets/d6f7dc5b-49d8-4ebd-889e-8bb0c8d66378)

## TODO

- [ ] Evaluation of SQL
  - ...
- [ ] Indices for faster lookups
  - ...
- [ ] Transactions
  - Atomicity and Durability: Transaction Log and WAL
  - Isolation: Locking
  - Consistency: Constraints and Validation
- [ ] Query Optimizer
  - ... idk
- [ ] Dealing with deadlocks
  - ... 
- [ ] Performance in storage w/ file system
  - Compression
  - ...
- [ ] Logging

## Reference
- https://github.com/awelm/simpledb
- https://github.com/AnarchistHoneybun/rdms
- https://github.com/CreatorsStack/CreatorDB
- https://github.com/jameycribbs/hare
- https://github.com/jameycribbs/ivy
