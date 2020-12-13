-- name: SelectNumberOfLikesByTodoID :one
SELECT COUNT(*)
FROM "todo_like"
WHERE "todo_id" = $1;
-- name: SelectAllUsersWhoLikeTodoByTodoID :many
SELECT *
FROM "user_account"
WHERE "user_id" IN (
        SELECT "owner_id"
        FROM "todo_like"
        WHERE "todo_id" = $1
    );
-- name: SelectLikeByOwnerIDAndTodoID :one
SELECT *
FROM "todo_like"
WHERE "owner_id" = $1
    AND "todo_id" = $2;
-- name: CreateTodoLike :exec
INSERT INTO "todo_like" ("owner_id", "todo_id")
VALUES ($1, $2)
ON CONFLICT ("owner_id", "todo_id") DO NOTHING;
-- name: DeleteTodoLike :exec
DELETE FROM "todo_like"
WHERE "owner_id" = $1
    AND "todo_id" = $2;