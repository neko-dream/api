-- name: CreateUserImage :exec
INSERT INTO user_images (
    user_images_id,
    user_id,
    image_url,
    width,
    height,
    extension,
    archived,
    url
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
