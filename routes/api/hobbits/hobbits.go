package hobbits

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func BuildHandleAPIPostHobbit(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recievedHobbit := models.Hobbit{}
		err := json.NewDecoder(r.Body).Decode(&recievedHobbit)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := models.User{}

		// TODO: Add error handling here
		err = db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sanitizedHobbit := models.Hobbit{
			Name:        strings.TrimSpace(recievedHobbit.Name),
			Image:       strings.TrimSpace(recievedHobbit.Image),
			Description: strings.TrimSpace(recievedHobbit.Description),
			User:        user,
		}

		err = db.Create(&sanitizedHobbit).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Infof("Created hobbit %+v", sanitizedHobbit)

		json.NewEncoder(w).Encode(sanitizedHobbit)
	}
}

func BuildHandleAPIGetHobbits(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hobbits := []models.Hobbit{}
		err := db.Joins("User").Find(&hobbits).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(hobbits)
	}
}

func BuildHandleAPIGetHobbit(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: add error handling
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		numericID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hobbit := models.Hobbit{}
		err = db.Where(models.Hobbit{ID: uint(numericID)}).Joins("User").First(&hobbit).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(hobbit)
	}
}

func BuildHandleAPIPutHobbit(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: add error handling
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		numericID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		currentHobbit := models.Hobbit{}
		err = db.Where(models.Hobbit{ID: uint(numericID)}).Joins("User").First(&currentHobbit).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Auth check: only the creator can update its hobbit
		if user.ID != currentHobbit.User.ID {
			log.Error("User does not match -> unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		recievedHobbit := models.Hobbit{}
		err = json.NewDecoder(r.Body).Decode(&recievedHobbit)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sanitizedHobbit := models.Hobbit{
			Name:        strings.TrimSpace(recievedHobbit.Name),
			Image:       strings.TrimSpace(recievedHobbit.Image),
			Description: strings.TrimSpace(recievedHobbit.Description),
		}
		db.Model(&currentHobbit).Updates(sanitizedHobbit)
		log.Infof("Updated hobbit %+v", sanitizedHobbit)
	}
}
