package authtocontext

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// AuthType is a alias for string used for AuthTypeSession and AuthTypePassPassword
// AuthType is used to store the type of used authentication within AuthDetails
type AuthType string

const (
	// AuthTypeSession is used to represent a session-authentication in AuthDetails
	AuthTypeSession AuthType = "Session"
	// AuthTypeAppPassword  used to represent a app-password-authentication in AuthDetails
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

// MiddlewareHandler is a middleware for handling all possible authentication options
// and storing them into the context of the http request
type MiddlewareHandler struct {
	next http.Handler
}

// New returns a new AuthToContextMiddlewareHandler
func New() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return MiddlewareHandler{
			next: next,
		}
	}
}

// ServeHTTP implements the core functionality of this middleware
func (m MiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := requestcontext.DB(r)
	log := requestcontext.Log(r)
	store := requestcontext.Store(r)

	handleBasicAuth := func(username string, password string) (AuthDetails, error) {
		if username == "" || password == "" {
			log.Warnf("username or password empty")
			return AuthDetails{}, fmt.Errorf("username or password empty")
		}

		user := &models.User{}
		if err := db.Where("username = ?", username).First(user).Error; err != nil {
			log.Warnf("found no user with username %s; %s ", username, err)
			return AuthDetails{}, err
		}

		appPasswords := []models.AppPassword{}

		if err := db.Joins("User").Where(models.AppPassword{UserID: user.ID}).Find(&appPasswords).Error; err != nil {
			log.Warnf("failed to find app passwords for user %s; %s ", username, err)
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
			log.Info("Password matched")
			if err := db.Model(&models.AppPassword{ID: matchedAppPassword.ID}).Updates(models.AppPassword{LastUsedAt: time.Now()}).Error; err != nil {
				log.Errorf("Failed to update LastUsedAt for app password %s", matchedAppPassword.ID)
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
			log.Debugf("Could not type assert cookie to AuthDetails, %+T", session.Values[AuthDetailsSessionKey])
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
			log.Debugln("User could not be authenticated by basicAuth")
		}
	} else if session, err := store.Get(r, "session"); err == nil && !session.IsNew && len(session.Values) > 0 {
		authDetails, err = handleSessionAuth(session)
		if err != nil {
			log.Debugln("User could not be authenticated by SessionAuth")
		}
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, AuthDetailsContextKey, authDetails)
	r = r.WithContext(ctx)

	m.next.ServeHTTP(w, r)
}
