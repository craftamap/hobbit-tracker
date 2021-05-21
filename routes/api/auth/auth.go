package auth

import (
	"encoding/json"
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func BuildHandleAPIGetAuth(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		contextAuthDetails := c.Value(authtocontext.AuthDetailsContextKey)
		authDetails, ok := contextAuthDetails.(authtocontext.AuthDetails)
		if !ok {
			json.NewEncoder(w).Encode(authtocontext.AuthDetails{
				Authenticated: false,
			})
			return
		}
		json.NewEncoder(w).Encode(authDetails)
	}
}
