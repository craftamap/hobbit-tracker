package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
)

func BuildHandleAPIGetAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		contextAuthDetails := c.Value(authtocontext.AuthDetailsContextKey)
		authDetails, ok := contextAuthDetails.(authtocontext.AuthDetails)
		if !ok {
			err := json.NewEncoder(w).Encode(authtocontext.AuthDetails{
				Authenticated: false,
			})
			if err != nil {
				slog.Error("failed to encode auth details", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		err := json.NewEncoder(w).Encode(authDetails)
		if err != nil {
			slog.Error("failed to encode auth details", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
