-- name: CreateToken :exec
insert into token (jti, user_id, kind, value, issued_at, expires_at)
values (?, ?, ?, ?, ?, ?);

-- name: GetValidToken :one
select * from token
where jti = ? and revoked = 0 and expires_at > CURRENT_TIMESTAMP;

-- name: GetTokenByJTI :one
select * from token where jti = ?;

-- name: GetTokensForUser :many
select * from token where user_id = ? order by issued_at desc;
