package api

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/routes/api/auth"
	"github.com/craftamap/hobbit-tracker/routes/api/hobbits"
	"github.com/craftamap/hobbit-tracker/routes/api/profile"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(api *mux.Router) {
	api.Use(func(h http.Handler) http.Handler {
		return handlers.ContentTypeHandler(h, "application/json")
	})

	rAuth := api.PathPrefix("/auth").Subrouter()
	auth.RegisterRoutes(rAuth)

	rHobbit := api.PathPrefix("/hobbits").Subrouter()
	hobbits.RegisterRoutes(rHobbit)
	rProfile := api.PathPrefix("/profile").Subrouter()
	profile.RegisterRoutes(rProfile)
}
