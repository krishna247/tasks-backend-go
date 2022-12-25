package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"io"
	"net/http"
	global "tasks/global"
	"tasks/model"
)

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createTask(w, r)
	case http.MethodGet:
		getTask(w, r)
	}

}

func createTask(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var task model.Task

	_ = json.Unmarshal(body, &task)
	fmt.Printf("%+v\n", task)
	taskId := insertTask(&task)
	fmt.Fprintf(w, taskId)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var getTaskInput model.GetTaskInput
	_ = json.Unmarshal(body, &getTaskInput)
	fmt.Printf("%+v\n", getTaskInput)
	tasks := getTasks(&getTaskInput)
	fmt.Printf("%+v\n", tasks)

	w.Header().Set("Content-Type", "application/json")
	tasksJson, _ := json.Marshal(tasks)
	w.Write(tasksJson)
}

func insertTask(task *model.Task) string {
	sqlStatement := `INSERT INTO task (user_uuid,deadline_date, priority, repeat_freq, tags, description, is_starred, is_done, sub_tasks)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	var taskId string
	global.DbConn.QueryRow(context.Background(), sqlStatement, task.UserUuid, task.DeadlineDate, task.Priority, task.RepeatFreq, task.Tags,
		task.Description, task.IsStarred, task.IsDone, task.SubTasks).Scan(&taskId)
	return taskId
}

func getTasks(getTaskInput *model.GetTaskInput) []model.Task {
	var tasks []model.Task
	pgxscan.Select(context.Background(), global.DbConn, &tasks, `SELECT * from task where user_uuid = $1`, getTaskInput.UserUuid)
	return tasks
}
