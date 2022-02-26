CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--
--  Create table user
--
CREATE TABLE "user_account" (
    "user_id" uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    "first_name" VARCHAR(256) NOT NULL,
    "last_name" VARCHAR(256) NOT NULL,
    "username" VARCHAR(256) NOT NULL UNIQUE,
    "email" VARCHAR(254) NOT NULL UNIQUE,
    "is_online" BOOLEAN NOT NULL DEFAULT FALSE,
    "last_login" TIMESTAMPTZ NULL,
    "password_hash" VARCHAR(128) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "user_account_username_f0fd97ac_like" ON "user_account" ("username" varchar_pattern_ops);
CREATE INDEX "user_account_email_d74bf2f6_like" ON "user_account" ("email" varchar_pattern_ops);
