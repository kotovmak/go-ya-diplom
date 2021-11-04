CREATE TABLE IF NOT EXISTS "user" (
  "user_id" bigserial PRIMARY KEY,
  "login" varchar NOT NULL,
  "password" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "withdrawn" bigint NOT NULL
);
CREATE TABLE IF NOT EXISTS "order" (
  "order_id" bigserial PRIMARY KEY,
  "number" varchar NOT NULL,
  "status" varchar NOT NULL,
  "accrual" varchar NOT NULL,
  "uploaded_at" timestamptz NOT NULL DEFAULT now(),
  "user_id" bigserial NOT NULL,
  CONSTRAINT "user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user" ("user_id")
);
CREATE TABLE IF NOT EXISTS "withdraw" (
  "withdraw_id" bigserial PRIMARY KEY,
  "order" varchar NOT NULL,
  "sum" bigint NOT NULL,
  "status" varchar NOT NULL,
  "processed_at" timestamptz NOT NULL DEFAULT now(),
  "user_id" bigserial NOT NULL,
  CONSTRAINT "user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user" ("user_id")
);