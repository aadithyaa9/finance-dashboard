# finance dashboard api

so this is a backend service for a role-based finance dashboard. handles user auth, permissions, and tracks income/expenses. 

built it primarily with go, postgres, and chi. 

## the stack
* **go 1.22** * **chi** - for routing. nice and lightweight, didn't want a massive framework.
* **postgresql** - using the numeric type for money stuff so we don't get floating point weirdness.
* **sqlx** - kept it simple. no ORM magic like gorm, just writing actual sql queries. 
* **jwt** - standard stateless auth.

## how to run this thing

you need go and postgres installed locally.

1. clone it and set up the env file:

   ```bash
   cp .env.example .env
   ```
2. set up the database. just run the migration script:
    
    ```bash
    createdb finance_db
    psql -d finance_db -f db/migrations/001_init_schema.sql
    ```
    
2. start it
   ```
   bash
    go mod tidy
    go run ./cmd
    ```
    
## How the roles work

there are 3 tiers of users:

- **viewer**: can basically only look at the dashboard and records.
    
- **analyst**: can do viewer stuff + create and update records.
    
- **admin**: can do everything, plus delete records and manage other users.
    

_note: right now, to get the first admin, you literally just pass `"role": "admin"` when hitting the `/api/auth/register` endpoint. obviously don't do this in prod, i need to write a db seeder for it._

## Thoughts & notes

- **why soft deletes?** financial records shouldn't really ever be hard deleted. so hitting the delete endpoint just slaps a `deleted_at` timestamp on it and hides it from the standard GET queries.
    
- **partial updates:** the update struct uses pointer fields (`*float64`, `*string`) so if you only want to update the notes on a record, you don't accidentally overwrite the amount to zero.
    
- **aggregations:** the dashboard stats (totals, trends, etc) are mostly handled directly in postgres using `COALESCE` and `SUM`. way faster than pulling everything into go memory to calculate it.