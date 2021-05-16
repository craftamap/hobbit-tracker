package main

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

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("jMcBBEBKAzw89XNb")
	store = sessions.NewCookieStore(key)
)

// AuthDetails is a struct used for storing authentication details within the cookie store
type AuthDetails struct {
	Authenticated bool   `json:"authenticated"`
	Username      string `json:"username,omitempty"`
	UserID        uint   `json:"userId,omitempty"`
}

// ContextKey is a string-alias used for safe-usage of the context method
type ContextKey string

const (
	// AuthDetailsContextKey is a key used to store and retrieve the authentication details to/from the http request
	AuthDetailsContextKey ContextKey = "AuthDetails"
	// AuthDetailsSessionKey is a key used to store and retrieve the authentication details to/from the http session
	AuthDetailsSessionKey string = "AuthDetails"
)

func init() {
	gob.Register(AuthDetails{})
}

type AuthToContextMiddlewareHandler struct {
	db   *gorm.DB
	log  *logrus.Logger
	next http.Handler
}

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
		return authDetails, nil
	}

	var authDetails AuthDetails
	var err error

	if username, password, ok := r.BasicAuth(); ok {
		authDetails, err = handleBasicAuth(username, password)
		if err != nil {
			m.log.Infof("User could not be authenticated by basicAuth")
		}
	} else if session, err := store.Get(r, "session"); err == nil && !session.IsNew {
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

func NewAuthToContextMiddlewareHandler(db *gorm.DB, log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return AuthToContextMiddlewareHandler{
			log:  log,
			db:   db,
			next: next,
		}
	}
}

type AuthMiddlewareHandler struct {
	log                   *logrus.Logger
	next                  http.Handler
	permitSessionAuth     bool
	permitAppPasswordAuth bool
}

func NewAuthMiddlewareHandler(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return AuthMiddlewareHandler{
			log:                   log,
			next:                  next,
			permitAppPasswordAuth: true,
			permitSessionAuth:     true,
		}
	}
}

func (m AuthMiddlewareHandler) PermitSessionAuth(permitSessionAuth bool) AuthMiddlewareHandler {
	m.permitSessionAuth = permitSessionAuth
	return m
}

func (m AuthMiddlewareHandler) PermitAppPasswordAuth(permitAppPasswordAuth bool) AuthMiddlewareHandler {
	m.permitAppPasswordAuth = permitAppPasswordAuth
	return m
}

// AuthMiddlewareBuilder creates a http.Handler which ensures that a user is authenticated.
// If a user is not authenticated, a http.StatusUnauthorized is written to the request and the request is returned.
// If a user is authenticated, the next handler is served.
func (m AuthMiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Basic auth for app authentication
	contextAuthDetails := r.Context().Value(AuthDetailsContextKey)
	authDetails := contextAuthDetails.(AuthDetails)
	if !authDetails.Authenticated {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	m.next.ServeHTTP(w, r)
}

// BuildHandleLogout is a function returning a http.HandlerFunc which logs out the current user.
// Users are getting logged out by setting their authDetails
func BuildHandleLogout(log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		authDetails := session.Values[AuthDetailsSessionKey].(AuthDetails)
		username := authDetails.Username

		session.Values[AuthDetailsSessionKey] = AuthDetails{
			Authenticated: false,
		}
		session.Save(r, w)

		redirectPath := r.PostForm.Get("redirect")
		if redirectPath == "" {
			redirectPath = "/"
		}

		w.Header().Add("Location", redirectPath)
		w.WriteHeader(http.StatusFound)

		log.Infof("Logged out user %s", username)
	}
}

// BuildHandleLogin is a function returning a http.HandlerFunc which logs in a user by their credentails.
func BuildHandleLogin(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Warn("Could not parse form data")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		username := r.PostForm.Get("username")
		if username == "" {
			log.Warn("request did not contain username")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		password := r.PostForm.Get("password")
		if password == "" {
			log.Warn("request did not contain password")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := &models.User{}
		if err := db.Where("username = ?", username).First(user).Error; err != nil {
			log.Warnf("found no user with username %s; %s ", username, err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user.Secret == "" {
			log.Warnf("found user with username %s, but no secret was found ", username)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Secret), []byte(password))
		if err != nil {
			log.Warnf("invalid password %s", username)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Auth successful
		session, _ := store.Get(r, "session")

		session.Values[AuthDetailsSessionKey] = AuthDetails{
			Authenticated: true,
			Username:      user.Username,
			UserID:        user.ID,
		}
		err = session.Save(r, w)
		log.Warnf("%s", err)

		redirectPath := r.PostForm.Get("redirect")
		if redirectPath == "" {
			redirectPath = "/"
		}

		w.Header().Add("Location", redirectPath)
		w.WriteHeader(http.StatusFound)

		log.Infof("Logged in user %s", username)
	}
}
