# Day 8 Notes: PostgreSQL, database/sql, and pgx

## 1. Connection Pooling (`*sql.DB`)

- `sql.Open("pgx", connStr)` does NOT open an immediate network connection; it initializes the pool manager handle. Use `db.PingContext(ctx)` to verify network reachability.
- **Pool Tuning Rules**:
  - `SetMaxOpenConns(N)`: Caps maximum active sockets to prevent Postgres socket starvation.
  - `SetMaxIdleConns(N)`: Keeps warm idle connections in reserve.
  - `SetConnMaxLifetime(D)`: Recycles old connections to prevent silent load balancer connection drops.

---

## 2. Querying & Row Scanning Mechanics

- **`ExecContext`**: Executing DDL/DML queries without returned rows (`CREATE`, `INSERT`, `UPDATE`).
- **`QueryRowContext`**: Fetching exactly one row (e.g., using `RETURNING` or `WHERE id = $1`).
- **`QueryContext`**: Fetching multi-row result sets.
- **`rows.Scan(&pointers...)`**: Decodes database wire types directly into memory addresses.
- **Resource Cleanup Checklist**:
  1. Always `defer rows.Close()` to return connections back to the pool.
  2. Always check `rows.Err()` after iteration loops to catch network streaming failures.

---

## 3. Atomic Transactions (`*sql.Tx`)

- `db.BeginTx(ctx, opts)` locks a single dedicated database socket connection to the transaction session.
- **The Rollback Pattern**:

  ```go
  tx, err := db.BeginTx(ctx, nil)
  if err != nil { return err }
  defer tx.Rollback() // Safe no-op if tx.Commit() runs first!

  // Execute queries on tx, NOT db
  if err := tx.Commit(); err != nil { return err }

  ```

- Use `SELECT ... FOR UPDATE` inside transactions to acquire row-level write locks against concurrent race conditions.

---
