-- name: CreateToken :one
insert into tokens (jti, user_id, kind, value, issued_at, expires_at)
values (?, ?, ?, ?, ?, ?)
returning *;

-- name: GetTokenByJTI :one
select * from tokens
where jti = ? and expires_at > current_timestamp;

-- name: GetTokensForUser :many
select * from tokens
where user_id = ? and expires_at > current_timestamp order by issued_at desc;
