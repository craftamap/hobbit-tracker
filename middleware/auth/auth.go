package auth

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/authToContext"
	"github.com/sirupsen/logrus"
)

type AuthMiddlewareHandlerBuilder struct {
	log                   *logrus.Logger
	permitSessionAuth     bool
	permitAppPasswordAuth bool
}

func Builder(log *logrus.Logger) AuthMiddlewareHandlerBuilder {
	return AuthMiddlewareHandlerBuilder{
		log: log,
	}
}

func (m AuthMiddlewareHandlerBuilder) WithPermitSessionAuth(permitSessionAuth bool) AuthMiddlewareHandlerBuilder {
	return AuthMiddlewareHandlerBuilder{
		log:                   m.log,
		permitAppPasswordAuth: m.permitAppPasswordAuth,
		permitSessionAuth:     permitSessionAuth,
	}
}

func (m AuthMiddlewareHandlerBuilder) WithPermitAppPasswordAuth(permitAppPasswordAuth bool) AuthMiddlewareHandlerBuilder {
	return AuthMiddlewareHandlerBuilder{
		log:                   m.log,
		permitAppPasswordAuth: permitAppPasswordAuth,
		permitSessionAuth:     m.permitSessionAuth,
	}
}

// AuthMiddlewareBuilder creates a http.Handler which ensures that a user is authenticated.
// If a user is not authenticated, a http.StatusUnauthorized is written to the request and the request is returned.
// If a user is authenticated, the next handler is served.
func (m AuthMiddlewareHandlerBuilder) Build(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextAuthDetails := r.Context().Value(authToContext.AuthDetailsContextKey)
		authDetails := contextAuthDetails.(authToContext.AuthDetails)

		if !authDetails.Authenticated {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if authDetails.AuthType == authToContext.AuthTypeAppPassword && !m.permitAppPasswordAuth {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("AuthType AppPassword not allowed for this endpoint"))
			return
		}

		if authDetails.AuthType == authToContext.AuthTypeSession && !m.permitSessionAuth {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("AuthType Session not allowed for this endpoint"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
