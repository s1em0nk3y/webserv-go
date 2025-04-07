ALTER TABLE Credentials DROP CONSTRAINT creds_uniq_user;
ALTER TABLE Referrals DROP CONSTRAINT referrals_uniq_id;
ALTER TABLE Actions DROP CONSTRAINT actions_uniq_name;