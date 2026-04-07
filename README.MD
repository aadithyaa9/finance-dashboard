# Finance Dashboard API

A robust backend service for a role-based financial management system. This API handles user authentication, multi-tier permissions, and tracks financial income and expenses with detailed analytics.

## Tech Stack

- **Language**: Go 1.26
    
- **Router**: Chi v5 (lightweight and idiomatic)
    
- **Database**: PostgreSQL (using `sqlx` for performance and `pgcrypto` for UUIDs)
    
- **Auth**: JWT (Stateless authentication)
    
- **Validation**: Go-Playground Validator v10
    

---

## Response Format

All API responses follow a consistent JSON envelope:

JSON

```
{
  "success": true,
  "data": { ... },     // Present on successful GET/POST
  "message": "...",    // Present on successful commands
  "error": "..."       // Present on failure
}
```

---

## API Endpoints

### 1. Authentication (Public)

|**Method**|**Endpoint**|**Description**|
|---|---|---|
|`POST`|`/api/auth/register`|Create a new account.|
|`POST`|`/api/auth/login`|Authenticate and receive a JWT.|

**Registration Payload:**

- `role`: Optional. Defaults to `viewer`. Accepted: `viewer`, `analyst`, `admin`.
    

---

### 2. Financial Records (Protected)

Requires a valid `Authorization: Bearer <token>` header.

|**Method**|**Endpoint**|**Min. Role**|**Description**|
|---|---|---|---|
|`GET`|`/api/records`|Viewer|List all records with filtering/pagination.|
|`GET`|`/api/records/{id}`|Viewer|Fetch details for a specific record.|
|`POST`|`/api/records`|Analyst|Create a new financial entry.|
|`PUT`|`/api/records/{id}`|Analyst|Update an existing record.|
|`DELETE`|`/api/records/{id}`|Admin|Soft-delete a record.|

**Filtering `/api/records`:**

Query parameters include `page`, `limit`, `type` (income/expense), `category`, `from` (YYYY-MM-DD), and `to`.

---

### 3. Dashboard Analytics (Protected)

Provides aggregated data for UI visualizations.

|**Method**|**Endpoint**|**Description**|
|---|---|---|
|`GET`|`/api/dashboard/summary`|Total income, expenses, and net balance.|
|`GET`|`/api/dashboard/by-category`|Spending/income totals grouped by category.|
|`GET`|`/api/dashboard/trends`|Monthly totals for the last 12 months.|
|`GET`|`/api/dashboard/recent`|The 10 most recent transactions.|

---

### 4. User Management (Admin Only)

|**Method**|**Endpoint**|**Description**|
|---|---|---|
|`GET`|`/api/users`|List all registered users.|
|`PATCH`|`/api/users/{id}/role`|Change a user's permission level.|
|`PATCH`|`/api/users/{id}/status`|Toggle `is_active` to enable/disable an account.|

---

## Role-Based Access Control (RBAC)

The system enforces three distinct permission levels:

1. **Viewer**: Read-only access to dashboard and personal records.
    
2. **Analyst**: All Viewer permissions + ability to **Create** and **Update** records.
    
3. **Admin**: Full system access, including **Deleting** records and managing **User** accounts.
    

---

## Local Setup

1. **Environment**: Copy `.env.example` to `.env` and configure your `DB_URL` and `JWT_SECRET`.
    
2. **Database**:
    
    Bash
    
    ```
    createdb finance_db
    psql -d finance_db -f db/migrations/001_init_schema.sql
    ```
    
3. **Run Server**:
    
    Bash
    
    ```
    go mod tidy
    go run ./cmd/main.go
    ```

## Thoughts & notes

- **why soft deletes?** financial records shouldn't really ever be hard deleted. so hitting the delete endpoint just slaps a `deleted_at` timestamp on it and hides it from the standard GET queries.
    
- **partial updates:** the update struct uses pointer fields (`*float64`, `*string`) so if you only want to update the notes on a record, you don't accidentally overwrite the amount to zero.
    
- **aggregations:** the dashboard stats (totals, trends, etc) are mostly handled directly in postgres using `COALESCE` and `SUM`. way faster than pulling everything into go memory to calculate it.
