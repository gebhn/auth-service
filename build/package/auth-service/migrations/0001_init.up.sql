create table if not exists users (
  user_id text primary key,
  username text unique not null,
  email text unique not null,
  password_hash text not null,
  created_at timestamp not null,
  updated_at timestamp not null
);

create table if not exists tokens (
  jti text primary key,
  user_id text not null references users(user_id) on delete cascade,
  kind text not null check (kind in ('access', 'refresh', 'password_reset', 'email_verification')),
  value text not null,
  issued_at timestamp not null,
  expires_at timestamp not null,
  created_at timestamp not null default current_timestamp
);

create table if not exists revoked_tokens (
  jti text primary key,
  revoked_at timestamp not null default current_timestamp
);

create index if not exists idx_token_user_id on tokens(user_id);
create index if not exists idx_token_expires_at on tokens(expires_at);
create index if not exists idx_token_kind on tokens(kind);

create index if not exists idx_user_email ON users(email);
create index if not exists idx_user_username ON users(username);

create trigger if not exists trigger_user_updated_at
after update on users
for each row
begin
    update user set updated_at = current_timestamp where user_id = old.user_id;
end;
