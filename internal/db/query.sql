-- name: GetDevice :one
SELECT * FROM devices
WHERE id = $1 LIMIT 1;

-- name: ListDevices :many
SELECT * FROM devices
ORDER BY id;

-- name: CreateDevice :one
INSERT INTO devices (
    name, token, last_sync, sleeps_until
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteDevice :exec
DELETE FROM devices
WHERE id = $1;

-- name: UpdateDevice :exec
UPDATE devices
SET name = $2, token = $3, last_sync = $4, sleeps_until = $5
WHERE id = $1;


-- name: GetImage :one
SELECT * FROM images
WHERE id = $1 LIMIT 1;

-- name: ListImages :many
SELECT * FROM images
ORDER BY id;

-- name: CreateImage :one
INSERT INTO images (
    device_id, permanent, data_original, data_processed
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteImage :exec
DELETE FROM images
WHERE id = $1;

-- name: UpdateImage :exec
UPDATE images
SET device_id = $2, permanent = $3, data_original = $4, data_processed = $5
WHERE id = $1;
