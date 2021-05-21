package authtocontext

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/craftamap/hobbit-tracker/models"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthType string

const (
	AuthTypeSession     AuthType = "Session"
	AuthTypeAppPassword AuthType = "AppPassword"
)

// ContextKey is a string-alias used for safe-usage of the context method
type ContextKey string

const (
	// AuthDetailsContextKey is a key used to store and retrieve the authentication details to/from the http request
	AuthDetailsContextKey ContextKey = "AuthDetails"
	// AuthDetailsSessionKey is a key used to store and retrieve the authentication details to/from the http session
	AuthDetailsSessionKey string = "AuthDetails"
)

// AuthDetails is a struct used for storing authentication details within the cookie store
type AuthDetails struct {
	Authenticated bool     `json:"authenticated"`
	Username      string   `json:"username,omitempty"`
	UserID        uint     `json:"userId,omitempty"`
	AuthType      AuthType `json:"authType,omitempty"`
}

func init() {
	gob.Register(AuthDetails{})
}

// AuthToContextMiddlewareHandler is a middleware for handling all possible authentication options
// and storing them into the context of the http request
type AuthToContextMiddlewareHandler struct {
	db    *gorm.DB
	log   *logrus.Logger
	store *sessions.CookieStore
	next  http.Handler
}

// New returns a new AuthToContextMiddlewareHandler
func New(db *gorm.DB, log *logrus.Logger, store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return AuthToContextMiddlewareHandler{
			log:   log,
			db:    db,
			next:  next,
			store: store,
		}
	}
}

// ServeHTTP implements the core functionality of this middleware
func (m AuthToContextMiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleBasicAuth := func(username string, password string) (AuthDetails, error) {
		user := &models.User{}
		if err := m.db.Where("username = ?", username).First(user).Error; err != nil {
			m.log.Warnf("found no user with username %s; %s ", username, err)
			return AuthDetails{}, err
		}

		appPasswords := []models.AppPassword{}

		if err := m.db.Joins("User").Where(models.AppPassword{UserID: user.ID}).Find(&appPasswords).Error; err != nil {
			m.log.Warnf("failed to find app passwords for user %s; %s ", username, err)
			return AuthDetails{}, err
		}
		secretAndPasswordMatched := false
		var matchedAppPassword models.AppPassword
		for _, appPassword := range appPasswords {
			err := bcrypt.CompareHashAndPassword([]byte(appPassword.Secret), []byte(password))
			if err == nil {
				secretAndPasswordMatched = true
				matchedAppPassword = appPassword
				break
			}
		}
		if secretAndPasswordMatched {
			m.log.Info("Password matched")
			if err := m.db.Model(&models.AppPassword{ID: matchedAppPassword.ID}).Updates(models.AppPassword{LastUsedAt: time.Now()}).Error; err != nil {
				m.log.Errorf("Failed to update LastUsedAt for app password %s", matchedAppPassword.ID)
				return AuthDetails{}, err
			}
			return AuthDetails{
				Authenticated: true,
				Username:      user.Username,
				UserID:        user.ID,
				AuthType:      AuthTypeAppPassword,
			}, nil
		}
		return AuthDetails{Authenticated: false}, fmt.Errorf("Could not find a matching app password")
	}

	handleSessionAuth := func(session *sessions.Session) (AuthDetails, error) {
		authDetails, ok := session.Values[AuthDetailsSessionKey].(AuthDetails)
		if !ok {
			m.log.Infof("Could not type assert cookie to AuthDetails, %+T", session.Values[AuthDetailsSessionKey])
			return AuthDetails{}, fmt.Errorf("Could not type assert cookie to AuthDetails, %+T", session.Values[AuthDetailsSessionKey])
		}
		if authDetails.AuthType != AuthTypeSession {
			authDetails.AuthType = AuthTypeSession
		}
		return authDetails, nil
	}

	var authDetails AuthDetails
	var err error

	if username, password, ok := r.BasicAuth(); ok {
		authDetails, err = handleBasicAuth(username, password)
		if err != nil {
			m.log.Infof("User could not be authenticated by basicAuth")
		}
	} else if session, err := m.store.Get(r, "session"); err == nil && !session.IsNew {
		authDetails, err = handleSessionAuth(session)
		if err != nil {
			m.log.Infof("User could not be authenticated by SessionAuth")
		}
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, AuthDetailsContextKey, authDetails)
	r = r.WithContext(ctx)

	m.next.ServeHTTP(w, r)
}
