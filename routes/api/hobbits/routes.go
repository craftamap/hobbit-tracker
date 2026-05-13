package hobbits

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/auth"
	"github.com/craftamap/hobbit-tracker/routes/api/records"
)

func GetRoutes() http.Handler {
	hobbitsRouter := http.NewServeMux()
	authMiddlewareBuilder := auth.Builder()

	hobbitsRouter.Handle("GET /api/hobbits/{$}", BuildHandleAPIGetHobbits())
	hobbitsRouter.Handle("POST /api/hobbits/{$}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPostHobbit()),
	))

	hobbitsRouter.Handle("GET /api/hobbits/{id}/{$}", BuildHandleAPIGetHobbit())
	hobbitsRouter.Handle("PUT /api/hobbits/{id}/{$}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIPutHobbit()),
	))
	hobbitsRouter.Handle("DELETE /api/hobbits/{id}/{$}", authMiddlewareBuilder.Build(
		http.HandlerFunc(BuildHandleAPIDeleteHobbit()),
	))

	hobbitsRouter.Handle("/api/hobbits/{hobbit_id}/records/", records.GetRoutes())

	return hobbitsRouter
}
