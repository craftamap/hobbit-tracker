package routes

import (
	"github.com/craftamap/hobbit-tracker/routes/api"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB, log *logrus.Logger, store sessions.Store) {

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", BuildHandleLogin()).Methods("POST")
	auth.HandleFunc("/logout", BuildHandleLogout())

	routeOfAPI := r.PathPrefix("/api").Subrouter()
	api.RegisterRoutes(routeOfAPI)
}
