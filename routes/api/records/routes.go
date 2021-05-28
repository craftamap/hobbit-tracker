package records

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/auth"
	"github.com/gorilla/mux"
)

func RegisterRoutes(records *mux.Router) {
	authMiddlewareBuilder := auth.Builder()

	records.Handle("/", BuildHandleAPIGetRecords()).Methods("GET")
	records.Handle("/", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPostRecord()),
	)).Methods("POST")
	records.Handle("/{record_id:[0-9]+}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPutRecord()),
	)).Methods("PUT")
	records.Handle("/{record_id:[0-9]+}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIDeleteRecord()),
	)).Methods("DELETE")
	records.Handle("/heatmap", BuildHandleAPIGetRecordsForHeatmap()).Methods("GET")
}
