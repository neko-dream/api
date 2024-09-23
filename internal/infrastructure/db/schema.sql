CREATE TYPE "AuthProvider" AS ENUM (
  'GOOGLE'
);

CREATE TYPE "session_status" AS ENUM (
  'ACTIVE',
  'INACTIVE'
);

CREATE TABLE "user_auths" (
  "user_id" uuid,
  "provider" AuthProvider,
  "subject" varchar,
  "created_at" timestamp
);

CREATE TABLE "sessions" (
  "session_id" uuid PRIMARY KEY,
  "user_id" uuid,
  "provider" AuthProvider,
  "session_status" session_status,
  "created_at" timestamp
);

CREATE TABLE "users" (
  "user_id" uuid PRIMARY KEY,
  "display_id" varchar,
  "display_name" varchar,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "talk_sessions" (
  "talk_session_id" uuid PRIMARY KEY,
  "theme" varchar NOT NULL,
  "finished_at" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "opinions" (
  "opinion_id" uuid PRIMARY KEY,
  "talk_session_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "opinionContent" varchar NOT NULL,
  "parent_opinion_id" uuid,
  "vote_id" uuid,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "votes" (
  "vote_id" uuid PRIMARY KEY,
  "opinion_id" uuid,
  "user_id" uuid,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "users_display_id_index" ON "users" ("display_id");

CREATE INDEX "talk_sessions_theme_index" ON "talk_sessions" ("theme");

ALTER TABLE "user_auths" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "opinions" ADD FOREIGN KEY ("talk_session_id") REFERENCES "talk_sessions" ("talk_session_id");

ALTER TABLE "opinions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "opinions" ADD FOREIGN KEY ("parent_opinion_id") REFERENCES "opinions" ("opinion_id");

ALTER TABLE "opinions" ADD FOREIGN KEY ("vote_id") REFERENCES "votes" ("vote_id");

ALTER TABLE "votes" ADD FOREIGN KEY ("opinion_id") REFERENCES "opinions" ("opinion_id");

ALTER TABLE "votes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

