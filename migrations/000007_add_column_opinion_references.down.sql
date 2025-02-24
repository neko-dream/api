-- opinions テーブルからPictureURL, ReferenceURL カラムを削除
ALTER TABLE "opinions" DROP COLUMN IF EXISTS "picture_url";
ALTER TABLE "opinions" DROP COLUMN IF EXISTS "reference_url";

