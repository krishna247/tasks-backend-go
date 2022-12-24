package main

import (
	"fmt"
	"log"
	"net/http"
	"tasks/handlers"
)

func main() {
	fmt.Println("Starting server")
	mux := http.NewServeMux()
	mux.HandleFunc("/user", handlers.UserHandler)
	mux.HandleFunc("/", handlers.HelloWorld)
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
