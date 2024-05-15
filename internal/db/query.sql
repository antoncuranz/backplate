-- name: GetDevice :one
SELECT * FROM devices
WHERE id = $1 LIMIT 1;

-- name: ListDevices :many
SELECT * FROM devices
ORDER BY id;

-- name: CreateDevice :one
INSERT INTO devices (
    name, token
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: DeleteDevice :exec
DELETE FROM devices
WHERE id = $1;

-- name: GetImage :one
SELECT * FROM images
WHERE id = $1 LIMIT 1;

-- name: ListImages :many
SELECT * FROM images
ORDER BY id;

-- name: CreateImage :one
INSERT INTO images (
    device_id, permanent
) VALUES (
    $1, $2
)
RETURNING *;

-- name: DeleteImage :exec
DELETE FROM images
WHERE id = $1;