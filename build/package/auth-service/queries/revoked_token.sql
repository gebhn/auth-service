-- name: CreateRevokedToken :one
insert into revoked_tokens (jti) values (?)
returning jti;

-- name: GetRevokedTokenByJti :one
select jti from revoked_tokens where jti = ? limit 1;

-- name: GetRevocableTokensByUser :many
SELECT jti
FROM tokens
WHERE user_id = ? AND kind = 'refresh'
ORDER BY issued_at DESC;
