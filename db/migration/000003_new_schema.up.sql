-- Active: 1710562887801@@127.0.0.1@5432@filesync
CREATE TABLE "file" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar,
  "user_id" integer,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "fileversion" (
  "id" SERIAL PRIMARY KEY,
  "file_id" integer,
  "version" integer,
  "updated_at" timestamp
);

CREATE TABLE "block" (
  "id" SERIAL PRIMARY KEY,
  "file_version_id" integer,
  "sequence" integer,
  "hash" VARCHAR
);

ALTER TABLE "fileversion" ADD FOREIGN KEY ("file_id") REFERENCES "file" ("id");

ALTER TABLE "block" ADD FOREIGN KEY ("file_version_id") REFERENCES "fileversion" ("id");