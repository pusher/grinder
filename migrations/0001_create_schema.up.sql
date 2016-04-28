CREATE TABLE users (
  user_id         UUID          PRIMARY KEY,
  user_name       VARCHAR(128)  NOT NULL,
  user_available  BOOLEAN       NOT NULL DEFAULT false,
  user_token      CHAR(64)
);
CREATE UNIQUE INDEX user_token_idx ON users(user_token);

CREATE TABLE matches (
  from_user  UUID     NOT NULL REFERENCES users(user_id),
  to_user    UUID     NOT NULL REFERENCES users(user_id),
  match      BOOLEAN  NOT NULL DEFAULT false
);
CREATE UNIQUE INDEX matches_users_idx ON matches (from_user, to_user);
