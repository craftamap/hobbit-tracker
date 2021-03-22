package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/craftamap/hobbit-tracker/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func BuildHandleApiGetAuth(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session")
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		authDetails, ok := session.Values["authDetails"].(AuthDetails)
		if !ok {
			json.NewEncoder(w).Encode(AuthDetails{
				Authenticated: false,
			})
			return
		}
		json.NewEncoder(w).Encode(authDetails)
	}
}

func BuildHandleApiPostHobbit(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
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
		err = db.Where("ID = ?", r.Context().Value(AUTH_DETAILS).(AuthDetails).UserID).First(&user).Error
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

func BuildHandleApiGetHobbits(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
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

func BuildHandleApiGetHobbit(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
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

func BuildHandleApiPutHobbit(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(AUTH_DETAILS).(AuthDetails).UserID).First(&user).Error
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

func BuildHandleApiPostRecord(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(AUTH_DETAILS).(AuthDetails).UserID).First(&user).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: add error handling
		vars := mux.Vars(r)
		hobbitId, ok := vars["hobbit_id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		numericHobbitId, err := strconv.ParseUint(hobbitId, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recievedRecord := models.NumericRecord{}
		err = json.NewDecoder(r.Body).Decode(&recievedRecord)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Are we allowed to create records for this hobbit?
		// TODO: error handling
		parentHobbit := models.Hobbit{}
		db.Where(models.Hobbit{ID: uint(numericHobbitId)}).Joins("User").First(&parentHobbit)

		if user.ID != parentHobbit.User.ID {
			log.Error("User does not match -> unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sanitizedRecord := models.NumericRecord{
			HobbitID: parentHobbit.ID,
			Value:    recievedRecord.Value,
			Comment:  recievedRecord.Comment,
		}

		err = db.Create(&sanitizedRecord).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		returnedRecord := sanitizedRecord
		// Do not put out hobbit id
		returnedRecord.HobbitID = 0

		json.NewEncoder(w).Encode(returnedRecord)
	}
}

func BuildHandleApiGetRecords(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hobbitId, ok := vars["hobbit_id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		numericHobbitId, err := strconv.ParseUint(hobbitId, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var records []models.NumericRecord
		err = db.Where(models.NumericRecord{
			HobbitID: uint(numericHobbitId),
		}).Find(&records).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(records)
	}
}
