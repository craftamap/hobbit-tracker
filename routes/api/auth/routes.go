package auth

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(auth *mux.Router, db *gorm.DB, log *logrus.Logger, store sessions.Store) {
	auth.HandleFunc("/", BuildHandleAPIGetAuth(db, log)).Methods("GET")
}
