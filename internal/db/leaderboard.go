package db

import "github.com/s1em0nk3y/webserv-go/internal/app"

func (db *DB) GetLeaderBoard() (app.LeaderBoard, error) {
	rows, err := db.db.Query(`
		SELECT
		  c.username user,
		  SUM(award) score
		FROM Tasklogs tl
		INNER JOIN Credentials c on c.id = tl.user_id
		GROUP BY c.username 
		ORDER BY score DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := app.LeaderBoard{}
	for rows.Next() {
		var name string
		var score float64
		err := rows.Scan(&name, &score)
		if err != nil {
			return nil, err
		}
		result = append(result, struct {
			User  string  "json:\"user\""
			Score float64 "json:\"score\""
		}{name, score},
		)
	}
	return result, nil
}
