package auth

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(auth *mux.Router) {
	auth.HandleFunc("/", BuildHandleAPIGetAuth()).Methods("GET")
}
