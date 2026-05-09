package routes

import (
	"log/slog"
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"golang.org/x/crypto/bcrypt"
)

// BuildHandleLogout is a function returning a http.HandlerFunc which logs out the current user.
// Users are getting logged out by setting their authDetails
func BuildHandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store := requestcontext.Store(r)

		session, _ := store.Get(r, "session")
		authDetails := session.Values[authtocontext.AuthDetailsSessionKey].(authtocontext.AuthDetails)
		username := authDetails.Username

		session.Values[authtocontext.AuthDetailsSessionKey] = authtocontext.AuthDetails{
			Authenticated: false,
		}
		// this deletes the session
		session.Options.MaxAge = -1
		err := session.Save(r, w)

		if err != nil {
			slog.Error("failed to save session", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		redirectPath := r.PostForm.Get("redirect")
		if redirectPath == "" {
			redirectPath = "/"
		}

		w.Header().Add("Location", redirectPath)
		w.WriteHeader(http.StatusFound)

		slog.Info("Logged out user", "username", username)
	}
}

// BuildHandleLogin is a function returning a http.HandlerFunc which logs in a user by their credentails.
func BuildHandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
		store := requestcontext.Store(r)

		err := r.ParseForm()
		if err != nil {
			slog.Warn("Could not parse form data")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		username := r.PostForm.Get("username")
		if username == "" {
			slog.Warn("request did not contain username")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		password := r.PostForm.Get("password")
		if password == "" {
			slog.Warn("request did not contain password")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := &models.User{}
		if err := db.Where("username = ?", username).First(user).Error; err != nil {
			slog.Warn("found no user with username", "username", username, "err", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user.Secret == "" {
			slog.Warn("found user with username, but no secret was found", "username", username)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Secret), []byte(password))
		if err != nil {
			slog.Warn("invalid password for user", "username", username)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Auth successful
		session, _ := store.Get(r, "session")

		session.Values[authtocontext.AuthDetailsSessionKey] = authtocontext.AuthDetails{
			Authenticated: true,
			Username:      user.Username,
			UserID:        user.ID,
		}
		err = session.Save(r, w)
		if err != nil {
			slog.Warn("failed to save session", "err", err)
		}

		redirectPath := r.PostForm.Get("redirect")
		if redirectPath == "" {
			redirectPath = "/"
		}

		w.Header().Add("Location", redirectPath)
		w.WriteHeader(http.StatusFound)

		slog.Info("Logged in user", "username", username)
	}
}
