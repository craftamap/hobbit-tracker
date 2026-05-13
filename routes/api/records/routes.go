package records

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/auth"
)

func GetRoutes() http.Handler {
	records := http.NewServeMux()
	authMiddlewareBuilder := auth.Builder()

	records.Handle("GET /api/hobbits/{hobbit_id}/records/{$}", BuildHandleAPIGetRecords())
	records.Handle("POST /api/hobbits/{hobbit_id}/records/{$}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPostRecord()),
	))
	records.Handle("PUT /api/hobbits/{hobbit_id}/records/{record_id}/{$}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPutRecord()),
	))
	records.Handle("DELETE /api/hobbits/{hobbit_id}/records/{record_id}/{$}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIDeleteRecord()),
	))

	return records
}
