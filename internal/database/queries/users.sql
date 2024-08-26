-- name: CheckUserExist :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE login = $1
) AS user_exist;

-- name: CreateUser :one
INSERT INTO users (login, password)
VALUES ($1, $2)
RETURNING id;

-- name: SaveRefreshToken :exec
UPDATE users
SET refresh_token = $2
WHERE id = $1;

-- name: GetUserByLogin :one
SELECT id, password
FROM users
WHERE login = $1;

-- name: Logout :exec
UPDATE users
SET refresh_token = ''
WHERE id = $1;

-- name: GetRefreshTokenById :one
SELECT refresh_token
FROM users
WHERE id = $1;