-- name: CreateUser :exec
insert into users (user_id, username, email, password_hash, created_at, updated_at)
values (?, ?, ?, ?, current_timestamp, current_timestamp);

-- name: UpdateUser :exec
UPDATE users
SET
  username = COALESCE(NULLIF(sqlc.arg(username), ''), username),
  email = COALESCE(NULLIF(sqlc.arg(email), ''), email),
  password_hash = COALESCE(NULLIF(sqlc.arg(password_hash), ''), password_hash)
WHERE
  user_id = sqlc.arg(user_id);

-- name: GetUserByID :one
select * from users where user_id = ?;

-- name: GetUserByEmail :one
select * from users where email = ?;

-- name: GetUserByUsername :one
select * from users where username = ?;
