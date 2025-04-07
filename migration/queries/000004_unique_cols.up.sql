ALTER TABLE Credentials ADD CONSTRAINT creds_uniq_user UNIQUE (username);
ALTER TABLE Referrals ADD CONSTRAINT referrals_uniq_id UNIQUE (id);
ALTER TABLE Actions ADD CONSTRAINT actions_uniq_name UNIQUE(name);