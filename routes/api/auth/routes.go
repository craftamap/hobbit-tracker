package auth

import "net/http"

func GetRoutes() http.Handler {
	authRoutes := http.NewServeMux()
	authRoutes.HandleFunc("GET /api/auth/", BuildHandleAPIGetAuth())
	return authRoutes
}
