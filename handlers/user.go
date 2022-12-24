package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUser(w, r)
	}

}

func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, uuid.New().String())
}
