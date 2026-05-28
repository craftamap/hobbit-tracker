package api

import (
	"fmt"
	"net/http"

	"github.com/craftamap/hobbit-tracker/httputil"
	"github.com/craftamap/hobbit-tracker/middleware/auth"
	authApi "github.com/craftamap/hobbit-tracker/routes/api/auth"
	"github.com/craftamap/hobbit-tracker/routes/api/hobbits"
	"github.com/craftamap/hobbit-tracker/routes/api/profile"
	"github.com/craftamap/hobbit-tracker/routes/api/records"
)

func GetRoutes() http.Handler {
	authMiddlewareBuilder := auth.Builder()

	apiRouter := http.NewServeMux()
	apiRouter.Handle("GET /api/auth", authApi.BuildHandleAPIGetAuth())
	apiRouter.Handle("GET /api/hobbits", hobbits.BuildHandleAPIGetHobbits())
	apiRouter.Handle("POST /api/hobbits", authMiddlewareBuilder.Build(hobbits.BuildHandleAPIPostHobbit()))

	apiRouter.Handle("GET /api/hobbits/{id}", hobbits.BuildHandleAPIGetHobbit())
	apiRouter.Handle("PUT /api/hobbits/{id}", authMiddlewareBuilder.Build(
		hobbits.BuildHandleAPIPutHobbit(),
	))
	apiRouter.Handle("DELETE /api/hobbits/{id}", authMiddlewareBuilder.Build(
		hobbits.BuildHandleAPIDeleteHobbit(),
	))
	apiRouter.Handle("GET /api/hobbits/{hobbit_id}/records", records.BuildHandleAPIGetRecords())
	apiRouter.Handle("POST /api/hobbits/{hobbit_id}/records", authMiddlewareBuilder.Build(
		records.BuildHandleAPIPostRecord(),
	))
	apiRouter.Handle("PUT /api/hobbits/{hobbit_id}/records/{record_id}", authMiddlewareBuilder.Build(
		records.BuildHandleAPIPutRecord(),
	))
	apiRouter.Handle("DELETE /api/hobbits/{hobbit_id}/records/{record_id}", authMiddlewareBuilder.Build(
		records.BuildHandleAPIDeleteRecord(),
	))
	apiRouter.Handle("GET /api/profile/me", authMiddlewareBuilder.Build(authApi.BuildHandleAPIGetAuth()))
	apiRouter.Handle("GET /api/profile/me/hobbits", authMiddlewareBuilder.Build(profile.BuildHandleProfileGetHobbits()))
	apiRouter.Handle("GET /api/profile/me/feed", authMiddlewareBuilder.Build(profile.GetMyFeed()))
	apiRouter.Handle("GET /api/profile/me/apppassword", authMiddlewareBuilder.Build(profile.BuildHandleGetAppPasswords()))
	apiRouter.Handle("POST /api/profile/me/apppassword", authMiddlewareBuilder.Build(profile.BuildHandlePostAppPassword()))
	apiRouter.Handle("DELETE /api/profile/me/apppassword/{id}", authMiddlewareBuilder.Build(profile.BuildHandleDeleteAppPassword()))
	apiRouter.Handle("GET /api/profile/{id}", authMiddlewareBuilder.Build(profile.GetOthersUserInfo()))
	apiRouter.Handle("GET /api/profile/{id}/follow", authMiddlewareBuilder.Build(profile.GetFollowForUser()))
	apiRouter.Handle("PUT /api/profile/{id}/follow", authMiddlewareBuilder.Build(profile.PutFollowForUser()))
	apiRouter.Handle("DELETE /api/profile/{id}/follow", authMiddlewareBuilder.Build(profile.DeleteFollowForUser()))
	apiRouter.Handle("GET /api/profile/{id}/hobbits", authMiddlewareBuilder.Build(profile.GetOthersHobbits()))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		apiRouter.ServeHTTP(w, r)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !(r.Method == http.MethodPut || r.Method == http.MethodPost || r.Method == http.MethodPatch) {
			handler.ServeHTTP(w, r)
			return
		}

		if httputil.IsContentType(r.Header.Get("Content-Type"), "application/json") {
			handler.ServeHTTP(w, r)
			return
		}

		http.Error(w, fmt.Sprintf("Unsupported content type %q; expected application/json", r.Header.Get("Content-Type")), http.StatusUnsupportedMediaType)
	})
}
