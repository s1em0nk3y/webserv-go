package db

import (
	"testing"

	"github.com/s1em0nk3y/webserv-go/internal/app"
)

func TestCompleteTask(t *testing.T) {
	user := "test1"
	messenger := "OTHER"
	action := "OTHER"
	taskDataStr := "somepost_id"

	taskData := app.TaskData{
		Messenger: &messenger,
		Action:    &action,
		TaskData:  &taskDataStr,
	}

	msg_id := -1
	if err := database.db.QueryRow("INSERT INTO Messengers(name) values($1) returning (id)", messenger).Scan(&msg_id); err != nil {
		t.Fatal(err)
	}
	defer database.db.Exec("DELETE FROM Messengers WHERE name = $1", messenger)
	act_id := -1
	if err := database.db.QueryRow("INSERT INTO Actions(name) values($1) returning (id)", action).Scan(&act_id); err != nil {
		t.Fatal(err)
	}
	defer database.db.Exec("DELETE FROM Actions WHERE name = $1", action)
	var id int
	if err := database.db.QueryRow(
		`INSERT INTO Credentials(username, passhash) values ($1, $2) returning (id)`,
		user, "test",
	).Scan(&id); err != nil {
		t.Fatal(err)
	}
	defer database.db.Exec("DELETE FROM Credentials where username = $1", user)

	task_id := -1
	if err := database.db.QueryRow(
		"INSERT INTO Tasks(messenger_id, action_id, task_data, award) values ($1, $2, $3, 100) RETURNING (id)",
		msg_id, act_id, taskDataStr,
	).Scan(&task_id); err != nil {
		t.Fatal(err)
	}
	defer database.db.Exec("DELETE FROM Tasks WHERE id = $1", task_id)
	award, err := database.CompleteUserTask(user, &taskData)
	if err != nil {
		t.Fatal(err)
	}
	defer database.db.Exec("DELETE FROM Tasklogs where user_id  = $1 and task_id = $2", id, task_id)
	if award == 0 {
		t.Fatal(err)
	}
}
