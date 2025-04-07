package db

import (
	"errors"

	_ "github.com/lib/pq"
)

func (db *DB) CreateUser(user string, passwordHash string) error {
	if _, err := db.db.Exec(
		"INSERT INTO Credentials (username, passhash) values ($1, $2)",
		user, passwordHash,
	); err != nil {
		db.logger.Err(err).Msg("an err occured when adding user")
		return err
	}
	return nil
}

func (db *DB) CheckCredents(user string, passwordHash string) (bool, error) {
	count := 0
	err := db.db.QueryRow(
		"SELECT COUNT(*) from Credentials where username = $1 and passhash = $2",
		user, passwordHash,
	).Scan(&count)
	if err != nil {
		db.logger.Err(err).Str("username", user).Msg("an err occured while checking credentials")
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (db *DB) ValidateUsername(username string) error {
	count := 0
	err := db.db.QueryRow(
		"SELECT COUNT(*) FROM Credentials where username = $1",
		username,
	).Scan(&count)
	if err != nil {
		db.logger.Err(err).Str("username", username).Msg("an err occured while checking username")
	}
	if count == 0 {
		return errors.New("user not found")
	}
	return nil
}
