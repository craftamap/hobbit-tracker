package auth

import (
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
)

0
0
0
0

// MiddlewareHandlerBuilder is a builder for authMiddlewareHandler, allowing easier configuration
type MiddlewareHandlerBuilder struct {
	permitSessionAuth     bool
	permitAppPasswordAuth bool
}

// Builder initializes the Builder with all required parameters of the builder
func Builder() MiddlewareHandlerBuilder {
	return MiddlewareHandlerBuilder{
		permitSessionAuth:     true,
		permitAppPasswordAuth: true,
	}
}

// WithPermitSessionAuth returns a new AuthMiddlewareHandlerBuilder with the new permitSessionAuth value
func (b MiddlewareHandlerBuilder) WithPermitSessionAuth(permitSessionAuth bool) MiddlewareHandlerBuilder {
	return MiddlewareHandlerBuilder{
		permitAppPasswordAuth: b.permitAppPasswordAuth,
		permitSessionAuth:     permitSessionAuth,
	}
}

// WithPermitAppPasswordAuth returns a new authMiddlewareHandlerBuilder with the permitAppPasswordAuth value
func (b MiddlewareHandlerBuilder) WithPermitAppPasswordAuth(permitAppPasswordAuth bool) MiddlewareHandlerBuilder {
	return MiddlewareHandlerBuilder{
		permitAppPasswordAuth: permitAppPasswordAuth,
		permitSessionAuth:     b.permitSessionAuth,
	}
}

// Build creates a new authMiddlewareHandler containing all the values of the Builder and the next handler given by parameter
func (b MiddlewareHandlerBuilder) Build(next http.Handler) http.Handler {
	return authMiddlewareHandler{
		permitSessionAuth:     b.permitSessionAuth,
		permitAppPasswordAuth: b.permitAppPasswordAuth,
		next:                  next,
	}
}

type authMiddlewareHandler struct {
	permitSessionAuth     bool
	permitAppPasswordAuth bool
	next                  http.Handler
}

// ServeHTTP implements the authentication
func (m authMiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := requestcontext.Log(r)

	contextAuthDetails := r.Context().Value(authtocontext.AuthDetailsContextKey)
	authDetails := contextAuthDetails.(authtocontext.AuthDetails)

	if !authDetails.Authenticated {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if authDetails.AuthType == authtocontext.AuthTypeAppPassword && !m.permitAppPasswordAuth {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("AuthType AppPassword not allowed for this endpoint"))
		if err != nil {
			log.Error(err)
		}
		return
	}

	if authDetails.AuthType == authtocontext.AuthTypeSession && !m.permitSessionAuth {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("AuthType Session not allowed for this endpoint"))
		if err != nil {
			log.Error(err)
		}
		return
	}

	m.next.ServeHTTP(w, r)
}
