package records

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(records *mux.Router, db *gorm.DB, log *logrus.Logger, store sessions.Store) {
	authMiddlewareBuilder := auth.Builder(log)

	records.Handle("/", BuildHandleAPIGetRecords(db, log)).Methods("GET")
	records.Handle("/", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPostRecord(db, log)),
	)).Methods("POST")
	records.Handle("/{record_id:[0-9]+}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPutRecord(db, log)),
	)).Methods("PUT")
	records.Handle("/{record_id:[0-9]+}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIDeleteRecord(db, log)),
	)).Methods("DELETE")
	records.Handle("/heatmap", BuildHandleAPIGetRecordsForHeatmap(db, log)).Methods("GET")

}
