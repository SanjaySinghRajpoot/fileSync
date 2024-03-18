CREATE TABLE "file" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar,
  "user_id" integer,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "fileversion" (
  "id" SERIAL PRIMARY KEY,
  "file_id" integer REFERENCES "file" ("id"),
  "version" integer,
  "updated_at" timestamp
);

CREATE TABLE "block" (
  "id" SERIAL PRIMARY KEY,
  "file_version_id" integer REFERENCES "fileversion" ("id"),
  "sequence" integer,
  "hash" VARCHAR
);
