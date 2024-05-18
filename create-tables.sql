DROP TABLE IF EXISTS records;
DROP TABLE IF EXISTS app_users;

CREATE TABLE app_users (
  id SERIAL   NOT NULL PRIMARY KEY,
  username    TEXT UNIQUE NOT NULL,
  password    TEXT NOT NULL,
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE records (
  id          SERIAL NOT NULL PRIMARY KEY,
  title       VARCHAR(64) NOT NULL,
  comment     VARCHAR(128),
  last_date   DATE NOT NULL,
  expiry_date DATE NOT NULL
);
