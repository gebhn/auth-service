create table if not exists user (
  user_id text primary key,
  username text unique not null,
  email text unique not null,
  password_hash text not null,
  created_at timestamp not null,
  updated_at timestamp not null
);

create table if not exists token (
  jti text primary key,
  user_id text not null references user(id) on delete cascade,
  kind text not null check (kind in ('access', 'refresh', 'password_reset', 'email_verification')),
  value text not null,
  issued_at timestamp not null,
  expires_at timestamp not null,
  created_at timestamp not null default current_timestamp
);

create table if not exists revoked_token (
  jti text primary key,
  revoked_at timestamp not null default current_timestamp
);

create index if not exists idx_token_user_id on token(user_id);
create index if not exists idx_token_expires_at on token(expires_at);
create index if not exists idx_token_kind on token(kind);

create index if not exists idx_user_email ON user(email);
create index if not exists idx_user_username ON user(username);

create trigger if not exists trigger_user_updated_at
after update on user
for each row
begin
    update user set updated_at = current_timestamp where user_id = old.user_id;
end;
