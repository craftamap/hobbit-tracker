package hobbits

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/auth"
	"github.com/craftamap/hobbit-tracker/routes/api/records"
	"github.com/gorilla/mux"
)

func RegisterRoutes(hobbits *mux.Router) {
	authMiddlewareBuilder := auth.Builder()

	hobbits.Handle("/", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPostHobbit()),
	)).Methods("POST")
	hobbits.Handle("/{id:[0-9]+}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPutHobbit()),
	)).Methods("PUT")
	hobbits.Handle("/{id:[0-9]+}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIDeleteHobbit()),
	)).Methods("DELETE")
	hobbits.Handle("/{id:[0-9]+}", BuildHandleAPIGetHobbit()).Methods("GET")
	hobbits.Handle("/", BuildHandleAPIGetHobbits()).Methods("GET")

	rRecords := hobbits.PathPrefix("/{hobbit_id:[0-9]+}/records").Subrouter()
	records.RegisterRoutes(rRecords)
}
