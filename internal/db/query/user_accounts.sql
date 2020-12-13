-- name: GetUserAccountByID :one
SELECT *
FROM user_account
WHERE user_id = $1;
-- name: GetUserAccountByEmail :one
SELECT *
FROM user_account
WHERE email = $1;
-- name: GetUserByEmailOrUsername :one
SELECT * FROM user_account
WHERE email = $1 OR username = $2;
-- name: GetAllUserAccounts :many
SELECT *
FROM user_account;
-- name: CreateUserAccount :one
INSERT INTO user_account(
        first_name,
        last_name,
        username,
        email,
        password_hash,
        is_online,
        created_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: UpdateUserAccountInfo :one
UPDATE user_account
SET first_name = $2,
    last_name = $3,
    email = $4
WHERE user_id = $1
RETURNING *;
-- name: SetUserPassword :one
UPDATE user_account
SET password_hash = $2
WHERE user_id = $1
RETURNING *;