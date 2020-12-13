--
-- Create model TodoLike
--
CREATE TABLE "todo_like" (
    "id" uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    "owner_id" uuid NOT NULL,
    "todo_id" uuid NOT NULL
);
ALTER TABLE "todo_like"
ADD CONSTRAINT "todo_like_owner_id_todo_id_0aca8b7a_uniq" UNIQUE ("owner_id", "todo_id");
ALTER TABLE "todo_like"
ADD CONSTRAINT "todo_like_owner_id_d17e997d_fk_user_account_user_id" FOREIGN KEY ("owner_id") REFERENCES "user_account" ("user_id") ON DELETE CASCADE;
ALTER TABLE "todo_like"
ADD CONSTRAINT "todo_like_todo_id_fbd44072_fk_todo_id" FOREIGN KEY ("todo_id") REFERENCES "todo" ("id") ON DELETE CASCADE;
CREATE INDEX "todo_like_owner_id_d17e997d" ON "todo_like" ("owner_id");
CREATE INDEX "todo_like_todo_id_fbd44072" ON "todo_like" ("todo_id");