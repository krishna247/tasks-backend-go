package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		global.LogError("Create task: Error in parsing message body", err)
		http.Error(w, "Error in parsing message body: "+err.Error(), 400)
		return
	}

	var task model.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		global.LogError("Create task: Error in parsing task json", err)
		http.Error(w, "Create task: Error in parsing task json: "+err.Error(), 400)
		return
	}

	validate := validator.New()
	err = validate.Struct(task)
	if err != nil {
		slog.Error("Error in validating task json", err)
		http.Error(w, "Error in validating task json: "+err.Error(), 400)
		return
	}

	taskId := insertTask(&task)
	taskResponse := map[string]string{"taskId": taskId}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(taskResponse)
	if err != nil {
		slog.Error("Error in generating response to task json", err)
		http.Error(w, "Error in generating response to task json: "+err.Error(), 500)
		return
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	var getTaskInput model.GetTaskInput
	err = json.Unmarshal(body, &getTaskInput)
	if err != nil {
		slog.Error("Get task: Error in parsing task json", err)
		http.Error(w, "Get task: Error in parsing task json: "+err.Error(), 400)
		return
	}

	fmt.Printf("%+v\n", getTaskInput)
	tasks := getTasks(&getTaskInput, &w)
	fmt.Printf("%+v\n", tasks)

	w.Header().Set("Content-Type", "application/json")
	tasksJson, _ := json.Marshal(tasks)
	_, err = w.Write(tasksJson)
	if err != nil {
		slog.Error("Get task: Error in sending response", err)
		http.Error(w, "Get task: Error in parsing task json: "+err.Error(), 400)
	}
}

func insertTask(task *model.Task) string {
	sqlStatement := `INSERT INTO task (user_uuid,deadline_date, priority, repeat_freq, tags, description, is_starred, is_done, sub_tasks)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	var taskId string
	err := global.DbConn.QueryRow(context.Background(), sqlStatement, task.UserUuid, task.DeadlineDate, task.Priority, task.RepeatFreq, task.Tags,
		task.Description, task.IsStarred, task.IsDone, task.SubTasks).Scan(&taskId)
	if err != nil {
		return ""
	}
	return taskId
}

func getTasks(getTaskInput *model.GetTaskInput, w *http.ResponseWriter) []model.Task {
	var tasks []model.Task
	err := pgxscan.Select(context.Background(), global.DbConn, &tasks, `SELECT * from task where user_uuid = $1`, getTaskInput.UserUuid)
	if err != nil {
		slog.Error("Get task: Error in getting tasks from db", err)
		http.Error(*w, "Get task: Error in getting tasks from db: "+err.Error(), 500)
	}
	return tasks
}
