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

# Some plan IDK

![image](https://github.com/user-attachments/assets/d6f7dc5b-49d8-4ebd-889e-8bb0c8d66378)

## TODO

- [ ] Evaluation of SQL
  - ...
- [ ] Indices for faster lookups
  - ...
- [ ] Transactions
  - ... ACID crap
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
