package hobbits

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/craftamap/hobbit-tracker/hub"
	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/gorilla/mux"
)

func BuildHandleAPIPostHobbit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
		eventHub := requestcontext.Hub(r)

		recievedHobbit := models.Hobbit{}
		err := json.NewDecoder(r.Body).Decode(&recievedHobbit)
		if err != nil {
			slog.Error("error encoding recieved hobbit", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := models.User{}

		// TODO: Add error handling here
		err = db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			slog.Error("error fetching user", "err", err)
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
			slog.Error("error creating hobbit", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		slog.Info("Created hobbit", "hobbit", sanitizedHobbit)

		err = json.NewEncoder(w).Encode(sanitizedHobbit)
		if err != nil {
			slog.Error("error encoding hobbit", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		eventHub.Broadcast(hub.ServerSideEvent{
			Typus:        hub.HobbitCreated,
			OptionalData: sanitizedHobbit,
		})
	}
}

func BuildHandleAPIGetHobbits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		hobbits := []models.Hobbit{}
		err := db.Joins("User").Find(&hobbits).Error
		if err != nil {
			slog.Error("failed to find hobbits", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(hobbits)
		if err != nil {
			slog.Error("failed to encode hobbits", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	}
}

func BuildHandleAPIGetHobbit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		// TODO: add error handling
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			slog.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		numericID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			slog.Error("failed to parse id", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hobbit := models.Hobbit{}
		err = db.Where(models.Hobbit{ID: uint(numericID)}).Joins("User").First(&hobbit).Error
		if err != nil {
			slog.Error("failed to find hobbit", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(w).Encode(hobbit)
		if err != nil {
			slog.Error("failed to encode hobbit", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func BuildHandleAPIPutHobbit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
		eventHub := requestcontext.Hub(r)

		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			slog.Error("failed to get user", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: add error handling
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			slog.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		numericID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			slog.Error("failed to parse id", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		currentHobbit := models.Hobbit{}
		err = db.Where(models.Hobbit{ID: uint(numericID)}).Joins("User").First(&currentHobbit).Error
		if err != nil {
			slog.Error("failed to find hobbit", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Auth check: only the creator can update its hobbit
		if user.ID != currentHobbit.User.ID {
			slog.Error("User does not match -> unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		recievedHobbit := models.Hobbit{}
		err = json.NewDecoder(r.Body).Decode(&recievedHobbit)
		if err != nil {
			slog.Error("error decoding recieved hobbit", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sanitizedHobbit := models.Hobbit{
			Name:        strings.TrimSpace(recievedHobbit.Name),
			Image:       strings.TrimSpace(recievedHobbit.Image),
			Description: strings.TrimSpace(recievedHobbit.Description),
		}
		db.Model(&currentHobbit).Updates(&sanitizedHobbit)
		// archivedAt is nullable and therefore isnt updatable with "Updates"
		db.Model(&currentHobbit).Update("archivedAt", recievedHobbit.ArchivedAt)
		slog.Info("Updated hobbit with values", "currentHobbit", currentHobbit, "values", sanitizedHobbit)

		err = json.NewEncoder(w).Encode(currentHobbit)
		if err != nil {
			slog.Error("error encoding hobbit", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		eventHub.Broadcast(hub.ServerSideEvent{
			Typus:        hub.HobbitModified,
			OptionalData: currentHobbit,
		})
	}
}

func BuildHandleAPIDeleteHobbit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
		eventHub := requestcontext.Hub(r)

		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			slog.Error("failed to get user", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: add error handling
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			slog.Error("Can't get id from mux")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		numericID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			slog.Error("failed to parse hobbit id", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		currentHobbit := models.Hobbit{}
		err = db.Where(models.Hobbit{ID: uint(numericID)}).Joins("User").First(&currentHobbit).Error
		if err != nil {
			slog.Error("failed to get hobbit", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Auth check: only the creator can update its hobbit
		if user.ID != currentHobbit.User.ID {
			slog.Error("User does not match -> unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = db.Model(&currentHobbit).Delete(&currentHobbit).Error
		if err != nil {
			slog.Error("failed to delete hobbit", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		slog.Info("Deleted hobbit", "hobbit", currentHobbit)

		eventHub.Broadcast(hub.ServerSideEvent{
			Typus:        hub.HobbitDeleted,
			OptionalData: currentHobbit,
		})
	}
}
