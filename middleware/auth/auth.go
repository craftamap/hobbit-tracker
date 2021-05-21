package auth

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/sirupsen/logrus"
)

// AuthMiddlewareHandlerBuilder is a builder for authMiddlewareHandler, allowing easier configuration
type AuthMiddlewareHandlerBuilder struct {
	log                   *logrus.Logger
	permitSessionAuth     bool
	permitAppPasswordAuth bool
}

// Builder initializes the Builder with all required parameters of the builder
func Builder(log *logrus.Logger) AuthMiddlewareHandlerBuilder {
	return AuthMiddlewareHandlerBuilder{
		log:                   log,
		permitSessionAuth:     true,
		permitAppPasswordAuth: true,
	}
}

// WithPermitSessionAuth returns a new AuthMiddlewareHandlerBuilder with the new permitSessionAuth value
func (b AuthMiddlewareHandlerBuilder) WithPermitSessionAuth(permitSessionAuth bool) AuthMiddlewareHandlerBuilder {
	return AuthMiddlewareHandlerBuilder{
		log:                   b.log,
		permitAppPasswordAuth: b.permitAppPasswordAuth,
		permitSessionAuth:     permitSessionAuth,
	}
}

// WithPermitAppPasswordAuth returns a new authMiddlewareHandlerBuilder with the permitAppPasswordAuth value
func (b AuthMiddlewareHandlerBuilder) WithPermitAppPasswordAuth(permitAppPasswordAuth bool) AuthMiddlewareHandlerBuilder {
	return AuthMiddlewareHandlerBuilder{
		log:                   b.log,
		permitAppPasswordAuth: permitAppPasswordAuth,
		permitSessionAuth:     b.permitSessionAuth,
	}
}

// Build creates a new authMiddlewareHandler containing all the values of the Builder and the next handler given by parameter
func (b AuthMiddlewareHandlerBuilder) Build(next http.Handler) http.Handler {
	return authMiddlewareHandler{
		log:                   b.log,
		permitSessionAuth:     b.permitSessionAuth,
		permitAppPasswordAuth: b.permitAppPasswordAuth,
		next:                  next,
	}
}

type authMiddlewareHandler struct {
	log                   *logrus.Logger
	permitSessionAuth     bool
	permitAppPasswordAuth bool
	next                  http.Handler
}

// ServeHTTP implements the authentication
func (m authMiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contextAuthDetails := r.Context().Value(authtocontext.AuthDetailsContextKey)
	authDetails := contextAuthDetails.(authtocontext.AuthDetails)

	if !authDetails.Authenticated {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if authDetails.AuthType == authtocontext.AuthTypeAppPassword && !m.permitAppPasswordAuth {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("AuthType AppPassword not allowed for this endpoint"))
		return
	}

	if authDetails.AuthType == authtocontext.AuthTypeSession && !m.permitSessionAuth {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("AuthType Session not allowed for this endpoint"))
		return
	}

	m.next.ServeHTTP(w, r)
}
