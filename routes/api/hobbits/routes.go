package hobbits

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/auth"
	"github.com/craftamap/hobbit-tracker/routes/api/records"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(hobbits *mux.Router, db *gorm.DB, log *logrus.Logger, store *sessions.CookieStore) {
	authMiddlewareBuilder := auth.Builder(log)

	hobbits.Handle("/", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPostHobbit(db, log)),
	)).Methods("POST")
	hobbits.Handle("/{id:[0-9]+}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPutHobbit(db, log)),
	)).Methods("PUT")
	hobbits.Handle("/{id:[0-9]+}", BuildHandleAPIGetHobbit(db, log)).Methods("GET")
	hobbits.Handle("/", BuildHandleAPIGetHobbits(db, log)).Methods("GET")

	rRecords := hobbits.PathPrefix("/{hobbit_id:[0-9]+}/records").Subrouter()
	records.RegisterRoutes(rRecords, db, log, store)
}
