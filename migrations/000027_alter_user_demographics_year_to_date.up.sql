-- year_of_birthカラムのデータを全て削除
UPDATE user_demographics SET year_of_birth = NULL;
-- カラム名をリネーム
ALTER TABLE user_demographics RENAME COLUMN year_of_birth TO date_of_birth;
