DROP TABLE IF EXISTS records;
DROP TABLE IF EXISTS app_users;

CREATE TABLE app_users (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username    TEXT UNIQUE NOT NULL,
  password    TEXT NOT NULL,
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE records (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title       VARCHAR(64) NOT NULL,
  comment     VARCHAR(128),
  last_date   DATE NOT NULL,
  created_by  UUID REFERENCES app_users(id) ON DELETE CASCADE
);
