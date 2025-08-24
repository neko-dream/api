-- ユーザー退会機能の削除

-- user_status_change_logsテーブルの削除
DROP TABLE IF EXISTS user_status_change_logs;

-- usersテーブルからwithdrawal_dateカラムを削除
ALTER TABLE users DROP COLUMN withdrawal_date;
