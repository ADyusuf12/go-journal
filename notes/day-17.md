# Day 17 Notes: Production Web + Database Integration (`pgx/v5`)

## 1. Clean Handler Dependency Injection (`AppContext`)

- Encapsulating shared resources (`*pgxpool.Pool`, `*slog.Logger`) inside an application context struct allows handlers to be defined as methods (`func (app *AppContext) HandlerName(...)`), avoiding global variables.
- Mirrors clean service/controller architectures without heavy framework boilerplate.

---

## 2. Request-Scoped Database Contexts & Timeouts

- Wrapping request contexts with timeouts (`context.WithTimeout(r.Context(), 3*time.Second)`) protects the database pool from hanging queries or unresponsive TCP socket states.
- Ensures that if a client disconnects or a query deadlocks, the goroutine and database connection are released instantly.

---

## 3. Parameterized Queries & `pgx` Execution

- Using positional parameters (`$1, $2, $3`) in `db.QueryRow()` completely prevents SQL injection vectors.
- Scanning returned primary keys (`.Scan(&insertedID)`) captures generated database IDs cleanly.
