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

func RegisterRoutes(profile *mux.Router, db *gorm.DB, log *logrus.Logger, store *sessions.CookieStore) {
	authMiddlewareBuilder := auth.Builder(log)

	profileMe := profile.PathPrefix("/me").Subrouter()
	profileMe.Use(authMiddlewareBuilder.Build)
	profileMe.Handle("/", apiAuth.BuildHandleAPIGetAuth(db, log))
	profileMe.Handle("/hobbits", http.HandlerFunc(BuildHandleAPIProfileGetHobbits(db, log))).Methods("GET")
	profileMeAppPassword := profileMe.PathPrefix("/apppassword").Subrouter()
	profileMeAppPassword.Use(authMiddlewareBuilder.WithPermitAppPasswordAuth(false).Build)
	profileMeAppPassword.HandleFunc("/", BuildHandleAPIProfileGetAppPasswords(db, log)).Methods("GET")
	profileMeAppPassword.HandleFunc("/", BuildHandleAPIProfilePostAppPassword(db, log)).Methods("POST")
}
