-- name: CreateUser :one
INSERT INTO users (id, username, email, hashed_password, created_at)
VALUES (
    gen_random_uuid ( ),
    $1,
    $2,
    $3,
    NOW()
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * from users 
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * from users 
WHERE username = $1;

-- name: GetUserByID :one
SELECT * from users 
WHERE id = $1;
