-- name: CreateUser :one
insert into users (user_id, username, email, password_hash, created_at, updated_at)
values (?, ?, ?, ?, current_timestamp, current_timestamp)
returning *;

-- name: GetUserByID :one
select * from users where user_id = ?;

-- name: GetUserByEmail :one
select * from users where email = ?;

-- name: GetUserByUsername :one
select * from users where username = ?;
