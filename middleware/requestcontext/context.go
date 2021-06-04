package requestcontext

import (
	"context"
	"net/http"

	"github.com/craftamap/hobbit-tracker/hub"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Handler implements a middleware allowing storage of relevant pointers for requests in the request context
type Handler struct {
	Store    sessions.Store
	DB       *gorm.DB
	Log      *logrus.Logger
	EventHub *hub.Hub
	next     http.Handler
}

// requestContextContextKey is a alias to string used for context keys
type requestContextContextKey string

const (
	requestContextStoreKey    requestContextContextKey = "store"
	requestContextDBKey       requestContextContextKey = "db"
	requestContextLogKey      requestContextContextKey = "log"
	requestContextEventHubKey requestContextContextKey = "eventHub"
)

// ServeHTTP implements the actual middleware code
func (m Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, requestContextStoreKey, m.Store)
	ctx = context.WithValue(ctx, requestContextDBKey, m.DB)
	ctx = context.WithValue(ctx, requestContextLogKey, m.Log)
	ctx = context.WithValue(ctx, requestContextEventHubKey, m.EventHub)

	r = r.WithContext(ctx)
	m.next.ServeHTTP(w, r)
}

// New creates a new Handler
func New(store sessions.Store, db *gorm.DB, log *logrus.Logger, eventHub *hub.Hub) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Handler{
			Store:    store,
			DB:       db,
			Log:      log,
			EventHub: eventHub,
			next:     next,
		}
	}
}

// Store retrieves the cookie store from the current request
func Store(r *http.Request) sessions.Store {
	return r.Context().Value(requestContextStoreKey).(sessions.Store)
}

// DB retrieves the database connection from the current request
func DB(r *http.Request) *gorm.DB {
	return r.Context().Value(requestContextDBKey).(*gorm.DB)
}

// Log retrieves the log from the current request
func Log(r *http.Request) *logrus.Logger {
	return r.Context().Value(requestContextLogKey).(*logrus.Logger)
}

// Hub retrieves the event hub from the current request
func Hub(r *http.Request) *hub.Hub {
	return r.Context().Value(requestContextEventHubKey).(*hub.Hub)
}
