package api

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/routes/api/auth"
	"github.com/gorilla/handlers"
)

func GetRoutes() http.Handler {
	apiRouter := http.NewServeMux()
	apiRouter.Handle("/auth/", http.StripPrefix("/auth", auth.GetRoutes()))
	// rAuth := api.PathPrefix("/auth").Subrouter()
	// auth.RegisterRoutes(rAuth)

	// rHobbit := api.PathPrefix("/hobbits").Subrouter()
	// hobbits.RegisterRoutes(rHobbit)
	// rProfile := api.PathPrefix("/profile").Subrouter()
	// profile.RegisterRoutes(rProfile)

	return handlers.ContentTypeHandler(apiRouter, "application/json")
}
