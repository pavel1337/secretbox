CREATE TABLE secrets (
  id            varchar(191)  NOT NULL PRIMARY KEY,
  content       text          NOT NULL,
  passphrase    text,
  expires_at    datetime(3)   NOT NULL
);

DROP TABLE secrets;