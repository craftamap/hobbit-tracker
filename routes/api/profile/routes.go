package profile

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/auth"
	apiAuth "github.com/craftamap/hobbit-tracker/routes/api/auth"
	"github.com/gorilla/mux"
)

func RegisterRoutes(profile *mux.Router) {
	authMiddlewareBuilder := auth.Builder()

	profileMe := profile.PathPrefix("/me").Subrouter()
	profileMe.Use(authMiddlewareBuilder.Build)
	profileMe.Handle("/", apiAuth.BuildHandleAPIGetAuth())
	profileMe.Handle("/hobbits", http.HandlerFunc(BuildHandleProfileGetHobbits())).Methods("GET")
	profileMeAppPassword := profileMe.PathPrefix("/apppassword").Subrouter()
	profileMeAppPassword.Use(authMiddlewareBuilder.Build) //.WithPermitAppPasswordAuth(false).Build)
	profileMeAppPassword.HandleFunc("/", BuildHandleGetAppPasswords()).Methods("GET")
	profileMeAppPassword.HandleFunc("/", BuildHandlePostAppPassword()).Methods("POST")
	profileMeAppPassword.HandleFunc("/{id:[0-9a-zA-Z\\-]+}", BuildHandleDeleteAppPassword()).Methods("DELETE")
}
