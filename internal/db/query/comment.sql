-- name: CountNumberOfCommentsByTodoID :one
SELECT COUNT(*)
FROM "comment"
WHERE "todo_id" = $1;
-- name: CreateComment :one
INSERT INTO "comment" (
        "owner_id",
        "content",
        "todo_id"
    )
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateChildCOmment :one
INSERT INTO "comment" (
        "owner_id",
        "content",
        "todo_id",
        "parent_comment_id"
    )
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: SelectMainCommentsByTodoId :many
SELECT *
FROM "comment"
WHERE "todo_id" = $1
    AND "parent_comment_id" IS NULL;
-- name: SelectSubcommentsByParentCommentId :many
SELECT *
FROM "comment"
WHERE "parent_comment_id" = $1;