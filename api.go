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

func BuildHandleAPIGetAuth(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		contextAuthDetails := c.Value(AuthDetailsContextKey)
		authDetails, ok := contextAuthDetails.(AuthDetails)
		if !ok {
			json.NewEncoder(w).Encode(AuthDetails{
				Authenticated: false,
			})
			return
		}
		json.NewEncoder(w).Encode(authDetails)
	}
}

func BuildHandleAPIProfileGetAppPasswords(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(AuthDetailsContextKey).(AuthDetails).UserID).First(&user).Error
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
		err = db.Where("ID = ?", r.Context().Value(AuthDetailsContextKey).(AuthDetails).UserID).First(&user).Error
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
		err := db.Where("ID = ?", r.Context().Value(AuthDetailsContextKey).(AuthDetails).UserID).First(&user).Error
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

func BuildHandleAPIPostRecord(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(AuthDetailsContextKey).(AuthDetails).UserID).First(&user).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: add error handling
		vars := mux.Vars(r)
		hobbitID, ok := vars["hobbit_id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		numericHobbitID, err := strconv.ParseUint(hobbitID, 10, 32)
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
		log.Info("recievedRecord", recievedRecord)

		// Are we allowed to create records for this hobbit?
		// TODO: error handling
		parentHobbit := models.Hobbit{}
		db.Where(models.Hobbit{ID: uint(numericHobbitID)}).Joins("User").First(&parentHobbit)

		if user.ID != parentHobbit.User.ID {
			log.Error("User does not match -> unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sanitizedRecord := models.NumericRecord{
			HobbitID:  parentHobbit.ID,
			Timestamp: recievedRecord.Timestamp,
			Value:     recievedRecord.Value,
			Comment:   recievedRecord.Comment,
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

func BuildHandleAPIPutRecord(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(AuthDetailsContextKey).(AuthDetails).UserID).First(&user).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: add error handling
		vars := mux.Vars(r)
		hobbitID, ok := vars["hobbit_id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recordID, ok := vars["record_id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		numericHobbitID, err := strconv.ParseUint(hobbitID, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		numericRecordID, err := strconv.ParseUint(recordID, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Are we allowed to create records for this hobbit?
		// TODO: error handling
		parentHobbit := models.Hobbit{}
		db.Where(models.Hobbit{ID: uint(numericHobbitID)}).Joins("User").First(&parentHobbit)

		if user.ID != parentHobbit.User.ID {
			log.Error("User does not match -> unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		recievedRecord := models.NumericRecord{}
		err = json.NewDecoder(r.Body).Decode(&recievedRecord)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Info("recievedRecord", recievedRecord)

		sanitizedRecord := models.NumericRecord{
			ID:        uint(numericRecordID),
			HobbitID:  parentHobbit.ID,
			Timestamp: recievedRecord.Timestamp,
			Value:     recievedRecord.Value,
			Comment:   recievedRecord.Comment,
		}
		// TODO: I think we should actually check first if the record even exists...
		err = db.Save(&sanitizedRecord).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		returnedRecord := sanitizedRecord

		json.NewEncoder(w).Encode(returnedRecord)

	}
}

func BuildHandleAPIDeleteRecord(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := db.Where("ID = ?", r.Context().Value(AuthDetailsContextKey).(AuthDetails).UserID).First(&user).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: add error handling
		vars := mux.Vars(r)
		hobbitID, ok := vars["hobbit_id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recordID, ok := vars["record_id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		numericHobbitID, err := strconv.ParseUint(hobbitID, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		numericRecordID, err := strconv.ParseUint(recordID, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Are we allowed to create records for this hobbit?
		// TODO: error handling
		parentHobbit := models.Hobbit{}
		db.Where(models.Hobbit{ID: uint(numericHobbitID)}).Joins("User").First(&parentHobbit)

		if user.ID != parentHobbit.User.ID {
			log.Error("User does not match -> unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Fetch Record first in order to properly return it as response
		deletedRecord := models.NumericRecord{
			ID: uint(numericRecordID),
		}

		err = db.Where(&deletedRecord).First(&deletedRecord).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recordToDelete := models.NumericRecord{
			ID: uint(numericRecordID),
		}
		err = db.Delete(&recordToDelete).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(deletedRecord)
	}
}

func BuildHandleAPIGetRecords(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hobbitID, ok := vars["hobbit_id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		numericHobbitID, err := strconv.ParseUint(hobbitID, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var records []models.NumericRecord
		err = db.Where(models.NumericRecord{
			HobbitID: uint(numericHobbitID),
		}).Find(&records).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(records)
	}
}

func BuildHandleAPIGetRecordsForHeatmap(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hobbitID, ok := vars["hobbit_id"]
		if !ok {
			log.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		numericHobbitID, err := strconv.ParseUint(hobbitID, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var records []models.NumericRecord
		err = db.Model(&models.NumericRecord{}).Select("hobbit_id, sum(value) as value, timestamp").Group("date(timestamp)").Where(models.NumericRecord{
			HobbitID: uint(numericHobbitID),
		}).Find(&records).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(records)
	}
}

func BuildHandleAPIProfileGetHobbits(db *gorm.DB, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hobbits := []models.Hobbit{}

		err := db.Joins("User").Where(&models.Hobbit{UserID: r.Context().Value(AuthDetailsContextKey).(AuthDetails).UserID}).Find(&hobbits).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(hobbits)
	}
}
