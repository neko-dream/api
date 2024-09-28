CREATE TABLE "user_auths" (
  "user_id" uuid,
  "provider" varchar NOT NULL,
  "subject" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
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

CREATE TABLE "users" (
  "user_id" uuid PRIMARY KEY,
  "display_id" varchar NOT NULL,
  "display_name" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "talk_sessions" (
  "talk_session_id" uuid PRIMARY KEY,
  "owner_id" uuid NOT NULL,
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
  "opinion_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "users_display_id_index" ON "users" ("display_id");

CREATE INDEX "talk_sessions_theme_index" ON "talk_sessions" ("theme");

ALTER TABLE "user_auths" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "talk_sessions" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("user_id");

ALTER TABLE "opinions" ADD FOREIGN KEY ("talk_session_id") REFERENCES "talk_sessions" ("talk_session_id");

ALTER TABLE "opinions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "opinions" ADD FOREIGN KEY ("parent_opinion_id") REFERENCES "opinions" ("opinion_id");

ALTER TABLE "opinions" ADD FOREIGN KEY ("vote_id") REFERENCES "votes" ("vote_id");

ALTER TABLE "votes" ADD FOREIGN KEY ("opinion_id") REFERENCES "opinions" ("opinion_id");

ALTER TABLE "votes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "users" ADD COLUMN "picture" varchar;
