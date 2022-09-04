package records

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/craftamap/hobbit-tracker/hub"
	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

func BuildHandleAPIPostRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.Tracer("hobbit").Start(r.Context(), "PostRecord")
		defer span.End()

		db := requestcontext.DB(r)
		log := requestcontext.Log(r)
		eventHub := requestcontext.Hub(r)

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

		err = json.NewEncoder(w).Encode(returnedRecord)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		eventHub.Broadcast(hub.ServerSideEvent{
			Typus:        hub.RecordCreated,
			OptionalData: sanitizedRecord,
		})
	}
}

func BuildHandleAPIPutRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
		log := requestcontext.Log(r)
		eventHub := requestcontext.Hub(r)

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

		err = json.NewEncoder(w).Encode(returnedRecord)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		eventHub.Broadcast(hub.ServerSideEvent{
			Typus:        hub.RecordModified,
			OptionalData: sanitizedRecord,
		})
	}
}

func BuildHandleAPIDeleteRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
		log := requestcontext.Log(r)
		eventHub := requestcontext.Hub(r)

		user := models.User{}

		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
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

		err = json.NewEncoder(w).Encode(deletedRecord)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		eventHub.Broadcast(hub.ServerSideEvent{
			Typus:        hub.RecordDeleted,
			OptionalData: deletedRecord,
		})
	}
}

func BuildHandleAPIGetRecords() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.Tracer("hobbit").Start(r.Context(), "GetRecords")
		defer span.End()

		db := requestcontext.DB(r)
		log := requestcontext.Log(r)

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
		span.AddEvent("DB")
		err = db.Where(models.NumericRecord{
			HobbitID: uint(numericHobbitID),
		}).Find(&records).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		span.AddEvent("DB done")

		err = json.NewEncoder(w).Encode(records)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func BuildHandleAPIGetRecordsForHeatmap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
		log := requestcontext.Log(r)

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

		err = json.NewEncoder(w).Encode(records)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
