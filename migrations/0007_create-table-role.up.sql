CREATE TABLE "role" (
    "id" SERIAL PRIMARY KEY,
    "code" VARCHAR(20) NOT NULL UNIQUE,
    "name" VARCHAR(20) NOT NULL UNIQUE
);
INSERT INTO "role" ("code", "name")
VALUES ('admin', 'Admin'),
    ('user', 'User');
--
-- Alter table user_account
--
ALTER TABLE "user_account"
ADD COLUMN "role_code" VARCHAR(20) NOT NULL REFERENCES "role" ("code") ON DELETE CASCADE DEFAULT 'user';