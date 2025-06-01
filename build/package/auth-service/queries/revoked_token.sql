-- name: RevokeToken :exec
insert into revoked_token (jti) values (?);

-- name: IsTokenRevoked :one
select exists (
  select 1 from revoked_token where jti = ?
);

-- name: RevokeAllRefreshTokensForUser :exec
insert into revoked_token (jti)
select jti from token
where user_id = ? and kind = 'refresh'
on conflict(jti) do nothing;
