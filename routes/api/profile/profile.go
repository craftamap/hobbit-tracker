package profile

import (
	"encoding/json"
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sethvargo/go-password/password"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func BuildHandleGetAppPasswords(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		appPasswords := []models.AppPassword{}
		if err = db.Where(&models.AppPassword{UserID: user.ID}).Find(&appPasswords).Error; err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sanitizedAppPasswords := []models.AppPassword{}
		for _, appPassword := range appPasswords {
			sanitizedAppPasswords = append(sanitizedAppPasswords, models.AppPassword{
				ID:          appPassword.ID,
				User:        appPassword.User,
				UserID:      appPassword.UserID,
				LastUsedAt:  appPassword.LastUsedAt,
				UpdatedAt:   appPassword.UpdatedAt,
				CreatedAt:   appPassword.CreatedAt,
				DeletedAt:   appPassword.DeletedAt,
				Description: appPassword.Description,
				Secret:      "", // DO NOT PRINT THE SECRET !!!
			})
		}
		json.NewEncoder(w).Encode(sanitizedAppPasswords)
	}
}

func BuildHandlePostAppPassword(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	// TODO: Limit total number of app passwords (10?)
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var dat map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&dat)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		descriptionI, ok := dat["description"]
		if !ok {
			log.Error("description missing in request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		description, ok := descriptionI.(string)
		if !ok || description == "" {
			log.Error("description missing in request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		generatedPassword, err := password.Generate(32, 10, 0, false, true)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// -1 = Use Default cost
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), -1)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		appPasswordToStore := models.AppPassword{
			Description: description,
			Secret:      string(encryptedPassword),
			UserID:      user.ID,
		}

		if err := db.Create(&appPasswordToStore).Error; err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sanitizedAppPassword := models.AppPassword{
			Secret:      generatedPassword,
			Description: appPasswordToStore.Description,
			UserID:      appPasswordToStore.UserID,
			ID:          appPasswordToStore.ID,
		}

		if err = json.NewEncoder(w).Encode(sanitizedAppPassword); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}

func BuildHandleDeleteAppPassword(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		parsedUuid, err := uuid.Parse(id)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		appPassword := models.AppPassword{
			ID: parsedUuid,
		}

		if err := db.First(&appPassword).Error; err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if appPassword.UserID != user.ID {
			log.Errorf(
				"Users for app password %s do not match! User %d is authenticated, but user %d is the owner of the app password",
				appPassword.ID, user.ID, appPassword.ID,
			)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := db.Delete(&appPassword).Error; err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		appPassword.Secret = "" // Even after creation, do not print secret

		json.NewEncoder(w).Encode(appPassword)
	}
}

func BuildHandleProfileGetHobbits(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hobbits := []models.Hobbit{}

		err := db.Joins("User").Where(&models.Hobbit{UserID: r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID}).Find(&hobbits).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(hobbits)
	}
}
