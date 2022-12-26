package main

import (
	"net/http"
	"tasks/db"
	"tasks/global"
	"tasks/handlers"
)

func init() {
	db.CreateDBConnection()
	global.SetupLog()
}

func main() {
	mux := http.NewServeMux()
	global.Log("Starting server")

	mux.HandleFunc("/task", handlers.TasksHandler)
	mux.HandleFunc("/user", handlers.UserHandler)
	mux.HandleFunc("/", handlers.HelloWorld)
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		global.LogError("Server Failure", err)
	}
}
