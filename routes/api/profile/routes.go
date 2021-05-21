package profile

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/auth"
	apiAuth "github.com/craftamap/hobbit-tracker/routes/api/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(profile *mux.Router, db *gorm.DB, log *logrus.Logger, store sessions.Store) {
	authMiddlewareBuilder := auth.Builder(log)

	profileMe := profile.PathPrefix("/me").Subrouter()
	profileMe.Use(authMiddlewareBuilder.Build)
	profileMe.Handle("/", apiAuth.BuildHandleAPIGetAuth(db, log))
	profileMe.Handle("/hobbits", http.HandlerFunc(BuildHandleProfileGetHobbits(db, log))).Methods("GET")
	profileMeAppPassword := profileMe.PathPrefix("/apppassword").Subrouter()
	profileMeAppPassword.Use(authMiddlewareBuilder.Build) //.WithPermitAppPasswordAuth(false).Build)
	profileMeAppPassword.HandleFunc("/", BuildHandleGetAppPasswords(db, log)).Methods("GET")
	profileMeAppPassword.HandleFunc("/", BuildHandlePostAppPassword(db, log)).Methods("POST")
	profileMeAppPassword.HandleFunc("/{id:[0-9a-zA-Z\\-]+}", BuildHandleDeleteAppPassword(db, log)).Methods("DELETE")
}
