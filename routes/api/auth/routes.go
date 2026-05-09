package auth

import "net/http"

func GetRoutes() http.Handler {
	authRoutes := http.NewServeMux()
	authRoutes.HandleFunc("GET /", BuildHandleAPIGetAuth())
	return authRoutes
}
