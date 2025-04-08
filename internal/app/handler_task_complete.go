package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type TaskCompleter interface {
	CompleteUserTask(username string, task *TaskData) (award float64, err error)
}

type TaskData struct {
	Messenger *string `json:"messenger"`
	Action    *string `json:"action"`
	TaskData  *string `json:"task_data"`
}

// POST /users/:id/task/complete
// should contain json like: {"messenger": "Telegram", "action": "like", "task_data": "some_post_id"}
func (a *App) completeTask(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["id"]
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	taskData := &TaskData{}
	if err := json.NewDecoder(r.Body).Decode(taskData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`please provide task data as: {"messenger": "Telegram", "action": "like", "task_data": "some_post_id"}`))
		return
	}
	award, err := a.TaskCompleter.CompleteUserTask(username, taskData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to complete task"))
		return
	}
	w.Write([]byte(fmt.Sprintf("%f", award)))
}
