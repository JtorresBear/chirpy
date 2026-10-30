-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(),NOW(),NOW(),$1,$2)
RETURNING  id, created_at, updated_at, email, is_chirpy_red ;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserPassword :one
SELECT * FROM users where email = $1;

-- name: UpdatePassword :one
UPDATE users SET hashed_password = $1, updated_at = NOW(), email = $2 WHERE id = $3
RETURNING id, created_at, updated_at, email, is_chirpy_red ;

-- name: UpgradToRed :exec
UPDATE users SET is_chirpy_red = TRUE, updated_at = NOW() WHERE id = $1;

-- name: GetUserById :one
SELECT id from users WHERE id = $1;