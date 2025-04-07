package db

import "errors"

func (db *DB) AddNewReferral(user string, referral string) error {
	if user == referral {
		return errors.New("user cant be his own referral")
	}
	_, err := db.db.Exec(`
		WITH userid as (
			SELECT id from Credentials
			WHERE username = $1
			LIMIT 1
		),
		referrerid as (
			SELECT id from Credentials
			WHERE username = $2
		)
		INSERT INTO Referrals (id, referrer_id)
		SELECT
			(SELECT id from userid),
			(SELECT id from referrerid)
		WHERE EXISTS (SELECT 1 from userid) AND EXISTS (SELECT 1 from referrerid)
	`, user, referral)
	db.logger.Err(err).Msg("adding new referral")
	return err
}
