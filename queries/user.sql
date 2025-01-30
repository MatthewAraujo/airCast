-- name: FinAllUsers :many
select * from user;

-- name: InsertUsers :one
INSERT INTO users (id, name, email, password)
VALUES (uuid_generate_v4(), $1, $2, $3)
RETURNING *;

-- name: FindUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: FindUserByEmail :one
SELECT * 
FROM users
WHERE email = $1
LIMIT 1;