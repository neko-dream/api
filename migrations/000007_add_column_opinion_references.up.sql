-- opinions テーブルにPictureURL, ReferenceURL カラムを追加
ALTER TABLE "opinions" ADD COLUMN "picture_url" varchar;
ALTER TABLE "opinions" ADD COLUMN "reference_url" varchar;
