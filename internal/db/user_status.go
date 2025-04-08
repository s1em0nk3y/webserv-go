package db

import "github.com/s1em0nk3y/webserv-go/internal/app"

func (db *DB) GetUserStatus(username string) (*app.UserStatus, error) {
	count := -1
	award := -1.0
	err := db.db.QueryRow(`
		SELECT
		  c.username,
		  COUNT (*),
		  SUM(award)
		FROM Tasklogs tl
		INNER JOIN Credentials c on c.id = tl.user_id
		WHERE c.username = $1
		GROUP BY c.username
	`, username,
	).Scan(&username, &count, &award)
	if err != nil {
		return nil, err
	}
	return &app.UserStatus{
		Username:           username,
		CompletedTaskCount: uint(count),
		Score:              award,
	}, nil
}
