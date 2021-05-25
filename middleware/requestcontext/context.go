package requestcontext

import (
	"context"
	"net/http"

	"github.com/craftamap/hobbit-tracker/hub"
	"github.com/sirupsen/logrus"
	"github.com/wader/gormstore/v2"
	"gorm.io/gorm"
)

type RequestContextMiddleware struct {
	Store    *gormstore.Store
	DB       *gorm.DB
	Log      *logrus.Logger
	EventHub *hub.Hub
	next     http.Handler
}

func (m RequestContextMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, "store", m.Store)
	ctx = context.WithValue(ctx, "db", m.DB)
	ctx = context.WithValue(ctx, "log", m.Log)
	ctx = context.WithValue(ctx, "eventHub", m.EventHub)

	r = r.WithContext(ctx)
	m.next.ServeHTTP(w, r)
}

func New(store *gormstore.Store, db *gorm.DB, log *logrus.Logger, eventHub *hub.Hub) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return RequestContextMiddleware{
			Store:    store,
			DB:       db,
			Log:      log,
			EventHub: eventHub,
			next:     next,
		}
	}
}

func Store(r *http.Request) *gormstore.Store {
	return r.Context().Value("store").(*gormstore.Store)
}

func DB(r *http.Request) *gorm.DB {
	return r.Context().Value("db").(*gorm.DB)
}

func Log(r *http.Request) *logrus.Logger {
	return r.Context().Value("log").(*logrus.Logger)
}

func Hub(r *http.Request) *hub.Hub {
	return r.Context().Value("eventHub").(*hub.Hub)
}
