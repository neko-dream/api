CREATE TABLE IF NOT EXISTS "user_group_info" (
    "talk_session_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "group_id" int NOT NULL,
    "pos_x" float NOT NULL,
    "pos_y" float NOT NULL,
    "updated_at" timestamp NOT NULL DEFAULT (now()),
    "created_at" timestamp NOT NULL DEFAULT (now()),

    PRIMARY KEY ("talk_session_id", "user_id")
);

CREATE TABLE IF NOT EXISTS "representative_opinions" (
    "talk_session_id" uuid NOT NULL,
    "opinion_id" uuid NOT NULL,
    "group_id" int NOT NULL,
    "rank" int NOT NULL,
    "updated_at" timestamp NOT NULL DEFAULT (now()),
    "created_at" timestamp NOT NULL DEFAULT (now()),

    PRIMARY KEY ("talk_session_id", "opinion_id", "group_id")
);

