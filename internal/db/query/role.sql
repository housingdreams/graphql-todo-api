
-- name: SelectRoleByCode :one
SELECT * FROM "role"
WHERE code = $1;