package auth

import (
	"encoding/json"
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/authToContext"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func BuildHandleAPIGetAuth(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		contextAuthDetails := c.Value(authToContext.AuthDetailsContextKey)
		authDetails, ok := contextAuthDetails.(authToContext.AuthDetails)
		if !ok {
			json.NewEncoder(w).Encode(authToContext.AuthDetails{
				Authenticated: false,
			})
			return
		}
		json.NewEncoder(w).Encode(authDetails)
	}
}
