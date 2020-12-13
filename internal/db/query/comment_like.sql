-- name: CountNumOfLikesOfCommentByID :one
SELECT COUNT(*)
FROM "comment_like"
WHERE "id" = $1;
-- name: SelectLikeByOwnerIDAndCommentID :one
SELECT *
FROM "comment_like"
WHERE "comment_id" = $1
    AND "owner_id" = $2;
-- name: LikeComment :exec
INSERT INTO "comment_like" ("owner_id", "comment_id")
VALUES ($1, $2) ON CONFLICT ("owner_id", "comment_id") DO NOTHING;
-- name: UnlikeComment :exec
DELETE FROM "comment_like"
WHERE "owner_id" = $1
    AND "comment_id" = $2;