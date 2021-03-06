package routes

import (
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
		log := requestcontext.Log(r)
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
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

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
func BuildHandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
		log := requestcontext.Log(r)
		store := requestcontext.Store(r)

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

		session.Values[authtocontext.AuthDetailsSessionKey] = authtocontext.AuthDetails{
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
