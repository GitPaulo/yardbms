# yardms
Yes A/nother Relational Database Management System.

```sh
yardms --storage=ram
```

In the REPL,

```sql
yardms> CREATE TABLE users;
Table users created
yardms> INSERT INTO users (id, name) VALUES (1, 'John Doe');
Row inserted into users
yardms> SELECT * FROM users;
Rows: [map[id:1 name:'John Doe']]
```

# Some plan IDK

![image](https://github.com/user-attachments/assets/d6f7dc5b-49d8-4ebd-889e-8bb0c8d66378)


## Reference
- https://github.com/awelm/simpledb
- https://github.com/AnarchistHoneybun/rdms
- https://github.com/CreatorsStack/CreatorDB
- https://github.com/jameycribbs/hare
- https://github.com/jameycribbs/ivy
