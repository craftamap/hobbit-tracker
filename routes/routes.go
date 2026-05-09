package routes

import (
	"net/http"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func GetRoutes(db *gorm.DB, store sessions.Store) http.Handler {
	//	auth := r.PathPrefix("/auth").Subrouter()
	//	auth.HandleFunc("/login", BuildHandleLogin()).Methods("POST")
	//	auth.HandleFunc("/logout", BuildHandleLogout())
	//
	//	routeOfAPI := r.PathPrefix("/api").Subrouter()
	//	api.RegisterRoutes(routeOfAPI)
	//
	return nil
}
