
DROP TABLE "user_auths" CASCADE;
DROP TABLE "users" CASCADE;
DROP TABLE "talk_session_locations" CASCADE;
DROP TABLE "talk_sessions" CASCADE;
DROP TABLE "opinions" CASCADE;
DROP TABLE "votes" CASCADE;
DROP TABLE "sessions" CASCADE;
DROP TABLE "user_demographics" CASCADE;

DROP EXTENSION IF EXISTS postgis CASCADE;
DROP TABLE IF EXISTS spatial_ref_sys CASCADE;
