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
			err := json.NewEncoder(w).Encode(authtocontext.AuthDetails{
				Authenticated: false,
			})
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		err := json.NewEncoder(w).Encode(authDetails)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
