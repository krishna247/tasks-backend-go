package handlers

import (
	"context"
	"encoding/json"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"strconv"
	global "tasks/global"
	"tasks/model"
)

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createTask(w, r)
	case http.MethodGet:
		getTask(w, r)
	case http.MethodDelete:
		deleteTask(w, r)
	default:
		http.Error(w, "Unsupported HTTP method", 400)

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

	taskId, err := insertTask(&task)
	if err != nil {
		slog.Error("Error in inserting task", err)
		http.Error(w, "Error in inserting task: "+err.Error(), 500)
		return
	}

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
	global.Log("Getting task for user: " + getTaskInput.UserUuid)
	tasks := getTasks(&getTaskInput, &w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	tasksJson, _ := json.Marshal(tasks)
	global.Log("Returning num tasks: " + strconv.Itoa(len(tasksJson)))

	_, err = w.Write(tasksJson)
	if err != nil {
		slog.Error("Get task: Error in sending response", err)
		http.Error(w, "Get task: Error in parsing task json: "+err.Error(), 400)
	}
}

func insertTask(task *model.Task) (string, error) {
	sqlStatement := `INSERT INTO task (user_uuid,deadline_date, priority, repeat_freq, tags, description, is_starred, is_done, sub_tasks)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	var taskId string
	err := global.DbConn.QueryRow(context.Background(), sqlStatement, task.UserUuid, task.DeadlineDate, task.Priority, task.RepeatFreq, task.Tags,
		task.Description, task.IsStarred, task.IsDone, task.SubTasks).Scan(&taskId)
	if err != nil {
		return "", err
	}
	return taskId, nil
}

func getTasks(getTaskInput *model.GetTaskInput, w *http.ResponseWriter) []model.Task {
	var tasks []model.Task
	err := pgxscan.Select(context.Background(), global.DbConn, &tasks, `SELECT * from task where user_uuid = '$1'`, getTaskInput.UserUuid)
	if err != nil {
		slog.Error("Get task: Error in getting tasks from db", err)
		http.Error(*w, "Get task: Error in getting tasks from db: "+err.Error(), 500)
	}
	return tasks
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	var taskInput model.DeleteTaskInput
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &taskInput)
	if err != nil {
		slog.Error("Delete task: Error in parsing request json", err)
		http.Error(w, "Delete task: Error in parsing request json: "+err.Error(), 400)
		return
	}

	_, err = global.DbConn.Exec(context.Background(), `delete from task where id=$1`, taskInput.TaskUuid)
	if err != nil {
		slog.Error("Delete task: Error in parsing request json", err)
		http.Error(w, "Delete task: Error in parsing request json: "+err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}
