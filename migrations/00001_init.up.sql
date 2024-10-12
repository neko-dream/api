CREATE EXTENSION IF NOT EXISTS postgis;
CREATE TABLE "users" (
  "user_id" uuid PRIMARY KEY,
  "display_id" varchar,
  "display_name" varchar,
  "icon_url" varchar,
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
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "user_demographics_user_id_index" ON "user_demographics" ("user_id");

CREATE TABLE "talk_sessions" (
  "talk_session_id" uuid PRIMARY KEY,
  "owner_id" uuid NOT NULL,
  "theme" varchar NOT NULL,
  "scheduled_end_time" timestamp NOT NULL,
  "finished_at" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "talk_session_locations" (
  "talk_session_id" uuid PRIMARY KEY,
  "location" GEOGRAPHY(POINT, 4326) NOT NULL,
  "city" varchar NOT NULL,
  "prefecture" varchar NOT NULL
);

CREATE TABLE "votes" (
  "vote_id" uuid PRIMARY KEY,
  "opinion_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "vote_type" SMALLINT NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "opinions" (
  "opinion_id" uuid PRIMARY KEY,
  "talk_session_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "parent_opinion_id" uuid, -- NULLならルート
  "title" varchar,
  "content" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

