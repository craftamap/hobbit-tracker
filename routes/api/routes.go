package api

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/routes/api/auth"
	"github.com/craftamap/hobbit-tracker/routes/api/hobbits"
	"github.com/craftamap/hobbit-tracker/routes/api/profile"
	"github.com/gorilla/handlers"
)

func GetRoutes() http.Handler {
	apiRouter := http.NewServeMux()
	apiRouter.Handle("/api/auth/", auth.GetRoutes())
	apiRouter.Handle("/api/hobbits/", hobbits.GetRoutes())
	apiRouter.Handle("/api/profile/", profile.GetRoutes())
	// rProfile := api.PathPrefix("/profile").Subrouter()
	// profile.RegisterRoutes(rProfile)

	return handlers.ContentTypeHandler(apiRouter, "application/json")
}
