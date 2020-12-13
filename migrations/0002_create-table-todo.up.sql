
CREATE TABLE "todo" (
    "id" uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    "title" varchar(256) NOT NULL UNIQUE,
    "content" text NOT NULL,
    "background" varchar(25) NOT NULL DEFAULT '#fff',
    "duedate" timestamp with time zone NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    "updated_at" timestamp with time zone NULL,
    "completed" boolean NOT NULL DEFAULT FALSE,
    "owner_id" uuid NOT NULL
);

ALTER TABLE "todo"
ADD CONSTRAINT "todo_owner_id_66a58107_fk_user_account_user_id" FOREIGN KEY ("owner_id") REFERENCES "user_account" ("user_id") ON DELETE CASCADE;
CREATE INDEX "todo_title_57d80d7a_like" ON "todo" ("title" varchar_pattern_ops);
CREATE INDEX "todo_completed_610887ca" ON "todo" ("completed");
CREATE INDEX "todo_owner_id_66a58107" ON "todo" ("owner_id");
