package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"tasks/global"
	"tasks/model"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "Unsupported HTTP method", 400)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		global.LogError("Create user: Error in parsing message body", err)
		http.Error(w, "Error in parsing message body: "+err.Error(), 400)
		return
	}

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		global.LogError("Create user: Error in parsing task json", err)
		http.Error(w, "Create user: Error in parsing task json: "+err.Error(), 400)
		return
	}

	var userUuid = uuid.New()
	sqlStatement := `INSERT INTO user_details (id, name, photo_url) VALUES ($1, $2, $3)`
	_, err = global.DbConn.Exec(context.Background(), sqlStatement, userUuid, user.Name, user.PhotoUrl)
	if err != nil {
		slog.Error("Get task: Error in getting tasks from db", err)
		http.Error(w, "Get task: Error in getting tasks from db: "+err.Error(), 500)
		return
	}

	var userResponse model.CreateUserResponse = model.CreateUserResponse{
		UserUuid: userUuid.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseJson, _ := json.Marshal(userResponse)
	_, err = w.Write(responseJson)
	if err != nil {
		slog.Error("Create user: Error in sending response", err)
		http.Error(w, "Create user: Error in parsing task json: "+err.Error(), 400)
	}
}
