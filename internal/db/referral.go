package db

import "errors"

func (db *DB) AddNewReferral(user string, referral string) error {
	if user == referral {
		return errors.New("user cant be his own referral")
	}
	id := -1
	err := db.db.QueryRow(`
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
		RETURNING id
		`, user, referral).Scan(&id)
	db.logger.Err(err).Msg("adding new referral")
	if err != nil {
		return err
	}
	if id == -1 {
		return errors.New("referral or user does not exist")
	}
	return nil
}
