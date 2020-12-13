--
-- Create model RefreshToken
--
CREATE TABLE "refresh_token" (
    "id" uuid NOT NULL PRIMARY KEY,
    "created_at" timestamp with time zone NOT NULL,
    "expires_at" timestamp with time zone NULL,
    "user_id" uuid NOT NULL
);

ALTER TABLE "refresh_token"
ADD CONSTRAINT "refresh_token_user_id_1d7a63ac_fk_user_account_user_id" FOREIGN KEY ("user_id") REFERENCES "user_account" ("user_id") ON DELETE CASCADE;
CREATE INDEX "refresh_token_user_id_1d7a63ac" ON "refresh_token" ("user_id");