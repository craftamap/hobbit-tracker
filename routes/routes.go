package routes

import (
	"github.com/craftamap/hobbit-tracker/middleware/auth"
	"github.com/craftamap/hobbit-tracker/routes/api"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB, log *logrus.Logger, store sessions.Store) {
	authMiddlewareBuilder := auth.Builder()
	r.Handle("/share", authMiddlewareBuilder.Build(HandleShare())).Methods("POST")
	r.Handle("/share/{id:[0-9]+}", authMiddlewareBuilder.Build(HandleGetShare())).Methods("GET")

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", BuildHandleLogin()).Methods("POST")
	auth.HandleFunc("/logout", BuildHandleLogout())

	routeOfAPI := r.PathPrefix("/api").Subrouter()
	api.RegisterRoutes(routeOfAPI)
}
