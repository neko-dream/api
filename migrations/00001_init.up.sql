CREATE TABLE "sessions" (
  "session_id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "provider" varchar NOT NULL,
  "session_status" int NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "last_activity_at" timestamp NOT NULL DEFAULT (now()),

  FOREIGN KEY ("user_id") REFERENCES "users" ("user_id")
);

CREATE TABLE "users" (
  "user_id" uuid PRIMARY KEY,
  "display_id" varchar,
  "display_name" varchar,
  "picture" varchar,
  "created_at" timestamp NOT NULL DEFAULT (now()),

  INDEX "users_display_id_index" ("display_id")
);


CREATE TABLE "user_auths" (
  "user_auth_id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "provider" varchar NOT NULL,
  "subject" varchar NOT NULL,
  "is_verified" boolean NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),

  FOREIGN KEY ("user_id") REFERENCES "users" ("user_id")
);

CREATE TABLE "user_demographics" (
  "user_id" uuid PRIMARY KEY,
  "year_of_birth" int NOT NULL,
  "occupation" tinyint NOT NULL,
  "gender" tinyint NOT NULL,
  "municiplaity" varchar NOT NULL,
  "household_size" tinyint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),

  FOREIGN KEY ("user_id") REFERENCES "users" ("user_id"),
  INDEX "user_demographics_year_of_birth_index" ("year_of_birth")
);

CREATE TABLE "talk_sessions" (
  "talk_session_id" uuid PRIMARY KEY,
  "owner_id" uuid NOT NULL,
  "theme" varchar NOT NULL,
  "finished_at" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now()),

  INDEX "talk_sessions_theme_index" ("theme"),

  FOREIGN KEY ("owner_id") REFERENCES "users" ("user_id")
);

CREATE TABLE "opinions" (
  "opinion_id" uuid PRIMARY KEY,
  "talk_session_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "opinionContent" varchar NOT NULL,
  "parent_opinion_id" uuid,
  "vote_id" uuid,
  "created_at" timestamp NOT NULL DEFAULT (now()),

  FOREIGN KEY ("talk_session_id") REFERENCES "talk_sessions" ("talk_session_id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("user_id"),
  FOREIGN KEY ("parent_opinion_id") REFERENCES "opinions" ("opinion_id"),
  FOREIGN KEY ("vote_id") REFERENCES "votes" ("vote_id")
);

CREATE TABLE "votes" (
  "vote_id" uuid PRIMARY KEY,
  "opinion_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),

  FOREIGN KEY ("user_id") REFERENCES "users" ("user_id"),
  FOREIGN KEY ("opinion_id") REFERENCES "opinions" ("opinion_id")
);
