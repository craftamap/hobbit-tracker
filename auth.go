package main

import (
	"context"
	"encoding/gob"
	"net/http"

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

type AuthDetails struct {
	Authenticated bool   `json:"authenticated"`
	Username      string `json:"username,omitempty"`
	UserID        uint   `json:"userId,omitempty"`
}

type ContextKey string

const (
	AUTH_DETAILS ContextKey = "authDetails"
)

func init() {
	gob.Register(AuthDetails{})
}

func AuthMiddlewareBuilder(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session")
			authDetails, ok := session.Values["authDetails"].(AuthDetails)
			if !ok {
				log.Infof("Could not type assert cookie to AuthDetails, %+T", session.Values["authDetails"])
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			auth := authDetails.Authenticated

			if !auth {
				log.Info("Could not type assert cookie to authDetails")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, AUTH_DETAILS, authDetails)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func BuildHandleLogout(log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		authDetails := session.Values["authDetails"].(AuthDetails)
		username := authDetails.Username

		session.Values["authDetails"] = AuthDetails{
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

		session.Values["authDetails"] = AuthDetails{
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
