package main

import (
	"net/http"
	"tasks/db"
	"tasks/handlers"
)

func init() {
	db.CreateDBConnection()
	SetupLog()
}

func main() {
	mux := http.NewServeMux()
	Log("Starting server")

	mux.HandleFunc("/task", handlers.TasksHandler)
	mux.HandleFunc("/user", handlers.UserHandler)
	mux.HandleFunc("/", handlers.HelloWorld)
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		LogError("Server Failure", err)
	}
}
