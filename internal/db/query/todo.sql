-- name: GetAllTodos :many
SELECT *
FROM "todo";
-- name: SelectAllTodosOfUserByUserID :many
SELECT *
FROM "todo"
WHERE "owner_id" = $1;
-- name: GetTodoByID :one
SELECT *
FROM "todo"
WHERE "id" = $1;
-- name: CreateTodo :one
INSERT INTO "todo" (
        "title",
        "content",
        "background",
        "duedate",
        "owner_id",
        "completed"
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: DeleteTodoByID :exec
DELETE FROM "todo"
WHERE "id" = $1
    AND "owner_id" = $2;
-- name: UpdateTodoByID :one
UPDATE "todo"
SET "title" = $2,
    "content" = $3,
    "background" = $4,
    "duedate" = $5,
    "updated_at" = $6,
    "completed" = $7
WHERE "id" = $1
    AND "owner_id" = $8
RETURNING *;