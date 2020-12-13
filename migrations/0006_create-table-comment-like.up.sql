--
-- Create model CommentLike
--
CREATE TABLE "comment_like" (
    "id" uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    "comment_id" uuid NOT NULL,
    "owner_id" uuid NOT NULL
);

ALTER TABLE "comment_like"
ADD CONSTRAINT "comment_like_owner_id_comment_id_7dbdfd2f_uniq" UNIQUE ("owner_id", "comment_id");
ALTER TABLE "comment_like"
ADD CONSTRAINT "comment_like_comment_id_ef80e084_fk_comment_id" FOREIGN KEY ("comment_id") REFERENCES "comment" ("id") ON DELETE CASCADE;
ALTER TABLE "comment_like"
ADD CONSTRAINT "comment_like_owner_id_03af57ba_fk_user_account_user_id" FOREIGN KEY ("owner_id") REFERENCES "user_account" ("user_id") ON DELETE CASCADE;
CREATE INDEX "comment_like_comment_id_ef80e084" ON "comment_like" ("comment_id");
CREATE INDEX "comment_like_owner_id_03af57ba" ON "comment_like" ("owner_id");