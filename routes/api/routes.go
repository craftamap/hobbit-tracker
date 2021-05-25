package api

import (
	"github.com/craftamap/hobbit-tracker/routes/api/auth"
	"github.com/craftamap/hobbit-tracker/routes/api/hobbits"
	"github.com/craftamap/hobbit-tracker/routes/api/profile"
	"github.com/gorilla/mux"
)

func RegisterRoutes(api *mux.Router) {
	rAuth := api.PathPrefix("/auth").Subrouter()
	auth.RegisterRoutes(rAuth)

	rHobbit := api.PathPrefix("/hobbits").Subrouter()
	hobbits.RegisterRoutes(rHobbit)
	rProfile := api.PathPrefix("/profile").Subrouter()
	profile.RegisterRoutes(rProfile)
}
