# Achieved

- ... 

# TODO

### General

- [ ] Logging
  - [ ] Implement basic logging for query execution and transactions
- [ ] Transactions
  - **Atomicity and Durability: Transaction Log and WAL**
    - [ ] Implement Write-Ahead Logging (WAL)
    - [ ] Ensure transaction logs are written before any changes
    - [ ] Add mechanisms for log recovery
  - **Isolation: Locking**
    - [ ] Implement row-level locking
    - [ ] Handle deadlock detection and resolution
  - **Consistency: Constraints and Validation**
    - [ ] Implement primary key constraints
    - [ ] Implement foreign key constraints
    - [ ] Implement unique constraints
    - [ ] Add support for check constraints

### Query Optimizer and Planner

- [ ] Query Optimizations
  - [ ] **Predicate Pushdown**
    - [ ] Implement basic predicate pushdown
    - [ ] Handle predicate pushdown in subqueries
  - [ ] **Projection Pushdown**
    - [ ] Identify and remove unused columns
    - [ ] Optimize SELECT statements by pushing down projections
  - [ ] **Join Order Optimization**
    - [ ] Implement heuristics for join order optimization
    - [ ] Use statistics to improve join order decisions
  - [ ] **Force Using Indexes**
    - [ ] Detect available indexes
    - [ ] Adjust query plans to utilize indexes
- [ ] Query Planner
  - [ ] **Cost-based Optimization**
    - [ ] Develop a cost model for query execution
    - [ ] Implement cost-based query plan selection
  - [ ] **Rule-based Optimization**
    - [ ] Define rules for query optimization
    - [ ] Apply rules to transform query plans
  - [ ] **Query Plan Caching**
    - [ ] Implement caching for query plans
    - [ ] Invalidate cache on schema changes

### SQL Engine

- [ ] Dealing with Deadlocks
  - [ ] Implement deadlock detection algorithm
  - [ ] Add deadlock resolution strategies (e.g., timeout, rollback)

### File Storage

- [ ] Don't Use JSON as File Storage lol
  - [ ] Choose a more efficient storage format (e.g., binary format, custom serialization)
- [ ] Performance in Storage w/ File System
  - [ ] Optimize file I/O operations
  - [ ] Implement data compression for storage
  - [ ] Design efficient data structures for storage
  - [ ] Handle large datasets with file partitioning
