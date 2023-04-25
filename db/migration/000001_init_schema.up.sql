CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" char(16)  NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT(now())
);
CREATE TABLE "users"(
  "id" char(16) PRIMARY KEY,
  "name"  varchar NOT NULL,
  "email" varchar NOT NULL ,
  "password"  varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT(now()),
  UNIQUE("email")
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("name");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users"("id") ON DELETE CASCADE;
