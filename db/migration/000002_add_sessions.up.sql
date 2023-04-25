CREATE TABLE "sessions"(
  "id" uuid PRIMARY KEY,
  "user_email"  varchar NOT NULL,
  "refresh_token" varchar NOT NULL ,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
alter table "sessions" add foreign key ("user_email") references "users"("email");