

CREATE TABLE "users" (
  "user_id" uuid PRIMARY KEY,
  "display_id" varchar,
  "display_name" varchar,
  "picture" varchar,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "session_id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "provider" varchar NOT NULL,
  "session_status" int NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "last_activity_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "user_auths" (
  "user_auth_id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "provider" varchar NOT NULL,
  "subject" varchar NOT NULL,
  "is_verified" boolean NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "user_demographics" (
  "user_demographics_id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "year_of_birth" int,
  "occupation" SMALLINT,
  "gender" SMALLINT NOT NULL,
  "municipality" varchar,
  "household_size" SMALLINT,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "talk_sessions" (
  "talk_session_id" uuid PRIMARY KEY,
  "owner_id" uuid NOT NULL,
  "theme" varchar NOT NULL,
  "finished_at" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "votes" (
  "vote_id" uuid PRIMARY KEY,
  "opinion_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
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

