CREATE TABLE user_images (
    user_images_id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    image_url TEXT NOT NULL,
    width INT NOT NULL,
    height INT NOT NULL,
    extension varchar(20) NOT NULL,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_user_images_user_id ON user_images(user_id, archived);
