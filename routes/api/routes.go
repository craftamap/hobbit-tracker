package api

import (
	"github.com/craftamap/hobbit-tracker/routes/api/auth"
	"github.com/craftamap/hobbit-tracker/routes/api/hobbits"
	"github.com/craftamap/hobbit-tracker/routes/api/profile"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(api *mux.Router, db *gorm.DB, log *logrus.Logger, store sessions.Store) {
	rAuth := api.PathPrefix("/auth").Subrouter()
	auth.RegisterRoutes(rAuth, db, log, store)

	rHobbit := api.PathPrefix("/hobbits").Subrouter()
	hobbits.RegisterRoutes(rHobbit, db, log, store)
	rProfile := api.PathPrefix("/profile").Subrouter()
	profile.RegisterRoutes(rProfile, db, log, store)
}
