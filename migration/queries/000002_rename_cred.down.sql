ALTER TABLE Credentials
  RENAME TO Users;
ALTER TABLE Users
  RENAME COLUMN username to nickname;