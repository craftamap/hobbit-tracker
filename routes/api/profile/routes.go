package profile

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/auth"
	apiAuth "github.com/craftamap/hobbit-tracker/routes/api/auth"
)

func GetRoutes() http.Handler {
	profile := http.NewServeMux()
	authMiddlewareBuilder := auth.Builder()

	profileMe := http.NewServeMux()
	profileMe.Handle("GET /api/profile/me", apiAuth.BuildHandleAPIGetAuth())
	profileMe.HandleFunc("GET /api/profile/me/hobbits", BuildHandleProfileGetHobbits())
	profileMe.HandleFunc("GET /api/profile/me/feed", GetMyFeed())

	profileMe.HandleFunc("GET /api/profile/me/apppassword", BuildHandleGetAppPasswords())
	profileMe.HandleFunc("POST /api/profile/me/apppassword/{$}", BuildHandlePostAppPassword())
	// todo: new id mapping
	profileMe.HandleFunc("DELETE /api/profile/me/apppassword/{id}", BuildHandleDeleteAppPassword())

	profile.Handle("/api/profile/me/", authMiddlewareBuilder.Build(profileMe))

	profileOthers := http.NewServeMux()
	// todo: new id mappings
	profileOthers.HandleFunc("GET /api/profile/{id}/", GetOthersUserInfo())
	profileOthers.HandleFunc("GET /api/profile/{id}/follow", GetFollowForUser())
	profileOthers.HandleFunc("PUT /api/profile/{id}/follow", PutFollowForUser())
	profileOthers.HandleFunc("DELETE /api/profile/{id}/follow", DeleteFollowForUser())
	profileOthers.HandleFunc("GET /api/profile/{id}/hobbits", GetOthersHobbits())
	profile.Handle("/api/profile/{id}/", authMiddlewareBuilder.Build(profileOthers))

	return profile
}
