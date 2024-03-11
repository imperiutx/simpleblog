CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "role" varchar NOT NULL DEFAULT 'user',
  "status" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "posts" (
  "id" bigserial PRIMARY KEY,
  "username" varchar,
  "title" varchar NOT NULL,
  "content" text NOT NULL,
  "tags" text[] NOT NULL,
  "status" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "edited_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "comments" (
  "id" bigserial PRIMARY KEY,
  "username" varchar,
  "post_id" bigint,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "edited_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "posts" ("username");

CREATE INDEX ON "comments" ("username");

CREATE INDEX ON "comments" ("post_id");

ALTER TABLE "posts" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "comments" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
