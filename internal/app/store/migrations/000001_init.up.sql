CREATE TABLE IF NOT EXISTS "users" (
  "user_id" bigserial PRIMARY KEY,
  "login" varchar NOT NULL,
  "password" varchar NOT NULL,
  "balance" bigint NOT NULL DEFAULT 0,
  "withdrawn" bigint NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS "orders" (
  "order_id" bigserial PRIMARY KEY,
  "number" varchar NOT NULL,
  "status" varchar NOT NULL,
  "accrual" bigint NOT NULL DEFAULT 0,
  "uploaded_at" timestamptz NOT NULL DEFAULT now(),
  "user_id" bigserial NOT NULL,
  CONSTRAINT "user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("user_id")
);
CREATE TABLE IF NOT EXISTS "withdraws" (
  "withdraw_id" bigserial PRIMARY KEY,
  "order" varchar NOT NULL,
  "sum" bigint NOT NULL,
  "status" varchar NOT NULL,
  "processed_at" timestamptz NOT NULL DEFAULT now(),
  "user_id" bigserial NOT NULL,
  CONSTRAINT "user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("user_id")
);
CREATE TABLE IF NOT EXISTS "query" (
  "query_id" bigserial PRIMARY KEY,
  "processing_at" timestamptz NOT NULL DEFAULT now(),
  "order" varchar NOT NULL
);