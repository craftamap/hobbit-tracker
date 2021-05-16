package profile

import (
	"encoding/json"
	"net/http"

	"github.com/craftamap/hobbit-tracker/middleware/authToContext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/sethvargo/go-password/password"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func BuildHandleAPIProfileGetAppPasswords(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(authToContext.AuthDetailsContextKey).(authToContext.AuthDetails).UserID).First(&user).Error
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
				Secret:      "",
			})
		}
		json.NewEncoder(w).Encode(sanitizedAppPasswords)
	}
}

func BuildHandleAPIProfilePostAppPassword(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	// TODO: Limit total number of app passwords (10?)
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := db.Where("ID = ?", r.Context().Value(authToContext.AuthDetailsContextKey).(authToContext.AuthDetails).UserID).First(&user).Error
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

func BuildHandleAPIProfileGetHobbits(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hobbits := []models.Hobbit{}

		err := db.Joins("User").Where(&models.Hobbit{UserID: r.Context().Value(authToContext.AuthDetailsContextKey).(authToContext.AuthDetails).UserID}).Find(&hobbits).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(hobbits)
	}
}
