CREATE TABLE talk_session_conclusions (
    talk_session_id uuid PRIMARY KEY,
    content TEXT NOT NULL,
    created_by uuid NOT NULL, -- 作成者
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_talk_session_conclusions_talk_session_id
    ON talk_session_conclusions(talk_session_id);

CREATE INDEX idx_talk_session_conclusions_creator
    ON talk_session_conclusions(created_by);

CREATE TABLE action_items (
    action_item_id uuid PRIMARY KEY,
    talk_session_id uuid NOT NULL,
    sequence INT NOT NULL, -- アクションアイテムの順番
    content TEXT NOT NULL, -- フローの内容
    status TEXT NOT NULL, -- 未着手, 進行中, 完了, 保留, 中止
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT check_status
        CHECK (status IN ('未着手', '進行中', '完了', '保留', '中止'))
);

