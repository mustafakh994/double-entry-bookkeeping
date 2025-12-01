-- name: CreateAccount :one
INSERT INTO accounts (balance, currency)
VALUES ($1, $2)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateAccountBalance :exec
UPDATE accounts
SET balance = $2
WHERE id = $1;

-- name: CreateTransaction :one
INSERT INTO transactions (from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListTransactions :many
SELECT * FROM transactions
WHERE from_account_id = $1 OR to_account_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
