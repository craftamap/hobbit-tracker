package routes

import (
	"github.com/craftamap/hobbit-tracker/routes/api"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", BuildHandleLogin()).Methods("POST")
	auth.HandleFunc("/logout", BuildHandleLogout())

	routeOfAPI := r.PathPrefix("/api").Subrouter()
	api.RegisterRoutes(routeOfAPI)
}
