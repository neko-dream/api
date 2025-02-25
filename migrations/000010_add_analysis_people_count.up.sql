ALTER TABLE representative_opinions ADD COLUMN "agree_count" INT NOT NULL DEFAULT 0;
ALTER TABLE representative_opinions ADD COLUMN "disagree_count" INT NOT NULL DEFAULT 0;
ALTER TABLE representative_opinions ADD COLUMN "pass_count" INT NOT NULL DEFAULT 0;
