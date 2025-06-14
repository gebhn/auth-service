-- name: CreateUser :exec
insert into users (user_id, username, email, password_hash, created_at, updated_at)
values (?, ?, ?, ?, current_timestamp, current_timestamp);

-- name: UpdateUser :exec
update users
set
  username = coalesce(nullif(sqlc.arg(username), ''), username),
  email = coalesce(nullif(sqlc.arg(email), ''), email),
  password_hash = coalesce(nullif(sqlc.arg(password_hash), ''), password_hash)
where
  user_id = sqlc.arg(user_id);

-- name: GetUserByID :one
select * from users where user_id = ?;

-- name: GetUserByEmail :one
select * from users where email = ?;

-- name: GetUserByUsername :one
select * from users where username = ?;
