--
-- Create model Comment
--
CREATE TABLE "comment" (
    "id" uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
    "updated_at" timestamp with time zone NULL,
    "owner_id" uuid NOT NULL,
    "parent_comment_id" uuid NULL,
    "todo_id" uuid NOT NULL
);

ALTER TABLE "comment"
ADD CONSTRAINT "comment_owner_id_a55b0a39_fk_user_account_user_id" FOREIGN KEY ("owner_id") REFERENCES "user_account" ("user_id") ON DELETE CASCADE;
ALTER TABLE "comment"
ADD CONSTRAINT "comment_parent_comment_id_e83b716b_fk_comment_id" FOREIGN KEY ("parent_comment_id") REFERENCES "comment" ("id") ON DELETE CASCADE;
ALTER TABLE "comment"
ADD CONSTRAINT "comment_todo_id_c80ba6b9_fk_todo_id" FOREIGN KEY ("todo_id") REFERENCES "todo" ("id") ON DELETE CASCADE;
CREATE INDEX "comment_owner_id_a55b0a39" ON "comment" ("owner_id");
CREATE INDEX "comment_parent_comment_id_e83b716b" ON "comment" ("parent_comment_id");
CREATE INDEX "comment_todo_id_c80ba6b9" ON "comment" ("todo_id");