CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "firstname" char NOT NULL,
  "lastname" char NOT NULL,
  "username" char NOT NULL,
  "type" char NOT NULL,
  "location" char NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);
CREATE TABLE "report" (
  "id" serial PRIMARY KEY,
  "user_id" int,
  "day" SMALLINT NOT NULL DEFAULT 0,
  "month" SMALLINT NOT NULL DEFAULT 0,
  "status" BOOLEAN DEFAULT false,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE TABLE "workdays" (
  "id" serial PRIMARY KEY,
  "user_id" bigint,
  "days" SMALLINT DEFAULT 0,
  "workhours" SMALLINT NOT NULL DEFAULT 0,
  "status" BOOLEAN DEFAULT false,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
  -- FOREIGN KEY (user_id) REFERENCES users (id)
);

