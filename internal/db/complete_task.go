package db

import "github.com/s1em0nk3y/webserv-go/internal/app"

func (db *DB) CompleteUserTask(username string, task *app.TaskData) (award float64, err error) {
	award = -1
	err = db.db.QueryRow(`
		WITH messengerid as (
			SELECT id from messengers WHERE name = $1 LIMIT 1
		),
		actionid as (
			SELECT id from actions WHERE name =  $2 LIMIT 1
		),
		userid as (
			SELECT id from Credentials
			WHERE username = $3
		),
		taskid as (
			SELECT id, award from Tasks
			WHERE messenger_id = (SELECT id FROM messengerid) AND
			  action_id = (SELECT id FROM actionid) AND task_data = $4
		),
		referralscount as (
			SELECT COUNT(*) as c FROM Referrals where referrer_id = ( SELECT id FROM userid )
		)
		INSERT INTO TaskLogs (user_id, task_id, award)
		SELECT
		  (SELECT id from userid),
		  (SELECT id from taskid),
		  (SELECT award from taskid) * ((SELECT c from referralscount) + 1)	
		RETURNING award
	`, *task.Messenger, *task.Action, username, *task.TaskData).Scan(&award)
	return award, err
}

/*

	WHERE
	  EXISTS (SELECT 1 from userid) AND
	  EXISTS (SELECT 1 from actionid) AND
	  EXISTS (SELECT 1 from messengerid) AND
	  EXISTS (SELECT 1 from taskid)

*/
