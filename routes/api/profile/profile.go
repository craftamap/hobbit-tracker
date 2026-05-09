package profile

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

func GetOthersUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		urlVariables := mux.Vars(r)
		otherUserStrId, ok := urlVariables["id"]
		if !ok {
			slog.Error("No user id found!")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		otherUserID, err := strconv.ParseUint(otherUserStrId, 10, 64)
		if err != nil {
			slog.Error("user id is not numeric", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// If userId is the userId of current user, redirect to /me
		//	authDetails := r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails)
		//	if authDetails.Authenticated && authDetails.UserID == uint(otherUserID) {
		//		http.Redirect(w, r, "../me", http.StatusTemporaryRedirect)
		//		return
		//	}

		otherUser := models.User{}
		if err = db.Where(models.User{ID: uint(otherUserID)}).First(&otherUser).Error; err != nil {
			slog.Error("user could not be found", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// sanitise otherUser for save output

		sanitisedOtherUser := models.User{
			ID:       otherUser.ID,
			Username: otherUser.Username,
			Image:    otherUser.Image,
		}

		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(sanitisedOtherUser)
		if err != nil {
			slog.Error("failed to encode response", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
func GetFollowForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		urlVariables := mux.Vars(r)
		otherUserStrId, ok := urlVariables["id"]
		if !ok {
			slog.Error("No user id found!")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		otherUserID, err := strconv.ParseUint(otherUserStrId, 10, 64)
		if err != nil {
			slog.Error("user id is not numeric", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// check if other user even exists

		otherUser := models.User{
			ID: uint(otherUserID),
		}

		if err := db.Where(&otherUser).First(&otherUser).Error; err != nil {
			slog.Error("failed to find user", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		authDetails := r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails)
		thisUserID := authDetails.UserID

		thisUser := models.User{
			ID: thisUserID,
		}
		if err := db.Preload("Follows").Where(&thisUser).First(&thisUser).Error; err != nil {
			slog.Error("failed to find user", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info("found user", "thisUser", thisUser)

		count := db.Model(&thisUser).Where(&otherUser).Association("Follows").Count()

		follows := count > 0

		json.NewEncoder(w).Encode(map[string]any{
			"follows": follows,
			"user":    otherUser,
		})
	}
}

func DeleteFollowForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		urlVariables := mux.Vars(r)
		otherUserStrId, ok := urlVariables["id"]
		if !ok {
			slog.Error("No user id found!")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		otherUserID, err := strconv.ParseUint(otherUserStrId, 10, 64)
		if err != nil {
			slog.Error("user id is not numeric", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// check if other user even exists

		otherUser := models.User{
			ID: uint(otherUserID),
		}

		if err := db.Where(&otherUser).First(&otherUser).Error; err != nil {
			slog.Error("failed to find user", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		authDetails := r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails)
		thisUserID := authDetails.UserID

		thisUser := models.User{
			ID: thisUserID,
		}
		if err := db.Preload("Follows").Where(&thisUser).First(&thisUser).Error; err != nil {
			slog.Error("failed to find user", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info("found user", "thisUser", thisUser)

		if err := db.Model(&thisUser).Association("Follows").Delete(&otherUser); err != nil {
			slog.Error("failed to unfollow user", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(otherUser)
	}
}
func PutFollowForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		urlVariables := mux.Vars(r)
		otherUserStrId, ok := urlVariables["id"]
		if !ok {
			slog.Error("No user id found!")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		otherUserID, err := strconv.ParseUint(otherUserStrId, 10, 64)
		if err != nil {
			slog.Error("user id is not numeric", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// check if other user even exists

		otherUser := models.User{
			ID: uint(otherUserID),
		}

		if err := db.Where(&otherUser).First(&otherUser).Error; err != nil {
			slog.Error("failed to find user", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authDetails := r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails)
		thisUserID := authDetails.UserID

		thisUser := models.User{
			ID: thisUserID,
		}
		if err := db.Preload("Follows").Where(&thisUser).First(&thisUser).Error; err != nil {
			slog.Error("failed to find user", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info("found user", "thisUser", thisUser)

		// check if we already follow the other user
		for _, follow := range thisUser.Follows {
			if follow.ID == uint(otherUserID) {
				slog.Error("User already follows user %d", "otherUserID", otherUserID)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if err := db.Model(&thisUser).Association("Follows").Append(&otherUser); err != nil {
			slog.Error("failed to follow user", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(otherUser)
	}
}

func GetOthersHobbits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		urlVariables := mux.Vars(r)
		otherUserStrId, ok := urlVariables["id"]
		if !ok {
			slog.Error("No user id found!")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		otherUserID, err := strconv.ParseUint(otherUserStrId, 10, 64)
		if err != nil {
			slog.Error("user id is not numeric", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// If userId is the userId of current user, redirect to /me
		// authDetails := r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails)
		// if authDetails.Authenticated && authDetails.UserID == uint(otherUserID) {
		// 	http.Redirect(w, r, "../me/hobbits", http.StatusTemporaryRedirect)
		// 	return
		// }

		hobbits := []models.Hobbit{}

		err = db.Joins("User").Where(&models.Hobbit{UserID: uint(otherUserID)}).Find(&hobbits).Error
		if err != nil {
			slog.Error("failed to find hobbits for other user", "err", err, "otherUserID", otherUserID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(hobbits)
		if err != nil {
			slog.Error("failed to encode response", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func BuildHandleGetAppPasswords() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		user := models.User{}

		// TODO: Add error handling here
		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			slog.Error("failed to find user", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		appPasswords := []models.AppPassword{}
		if err = db.Where(&models.AppPassword{UserID: user.ID}).Find(&appPasswords).Error; err != nil {
			slog.Error("failed to find app passwords", "err", err)
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
		err = json.NewEncoder(w).Encode(sanitizedAppPasswords)
		if err != nil {
			slog.Error("failed to encode response", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func BuildHandlePostAppPassword() http.HandlerFunc {
	// TODO: Limit total number of app passwords (10?)
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		user := models.User{}

		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			slog.Error("failed to find user", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var dat map[string]any
		err = json.NewDecoder(r.Body).Decode(&dat)
		if err != nil {
			slog.Error("failed to decode request body", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		descriptionI, ok := dat["description"]
		if !ok {
			slog.Error("description missing in request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		description, ok := descriptionI.(string)
		if !ok || description == "" {
			slog.Error("description missing in request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		generatedPassword, err := password.Generate(32, 10, 0, false, true)
		if err != nil {
			slog.Error("failed to generate password", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// -1 = Use Default cost
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), -1)
		if err != nil {
			slog.Error("failed to hash password", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		appPasswordToStore := models.AppPassword{
			Description: description,
			Secret:      string(encryptedPassword),
			UserID:      user.ID,
		}

		if err := db.Create(&appPasswordToStore).Error; err != nil {
			slog.Error("failed to create app password", "err", err)
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
			slog.Error("failed to encode response", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}

func BuildHandleDeleteAppPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		user := models.User{}

		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			slog.Error("failed to find user", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			slog.Error("No app password id found!")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			slog.Error("failed to parse app password id", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		appPassword := models.AppPassword{
			ID: parsedUUID,
		}

		if err := db.First(&appPassword).Error; err != nil {
			slog.Error("failed to find app password", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if appPassword.UserID != user.ID {
			slog.Error(
				"Users for app password do not match! User is authenticated, but user is the owner of the app password",
				"appPasswordId", appPassword.ID, "userId", user.ID, "appPasswordUserId", appPassword.UserID,
			)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := db.Delete(&appPassword).Error; err != nil {
			slog.Error("failed to delete app password", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		appPassword.Secret = "" // Even after creation, do not print secret

		err = json.NewEncoder(w).Encode(appPassword)
		if err != nil {
			slog.Error("failed to encode response", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func BuildHandleProfileGetHobbits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)

		hobbits := []models.Hobbit{}

		err := db.Joins("User").Where(&models.Hobbit{UserID: r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID}).Find(&hobbits).Error
		if err != nil {
			slog.Error("failed to find hobbits for user", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(hobbits)
		if err != nil {
			slog.Error("failed to encode response", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

type FeedEventTypus string

const (
	FeedEventTypusHobbitCreated FeedEventTypus = "HobbitCreated"
	FeedEventTypusRecordCreated FeedEventTypus = "RecordCreated"
)

type FeedEvent struct {
	FeedEventTypus FeedEventTypus
	CreatedAt      time.Time
	Payload        any
}

func GetMyFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
		authDetails := authtocontext.Get(r)

		user := models.User{}

		err := db.Preload("Follows").Where("ID = ?", authDetails.UserID).First(&user).Error
		if err != nil {
			slog.Error("failed to find user", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userIdsOfFollows := []uint{user.ID} // the user also "follows" itself
		for _, follow := range user.Follows {
			userIdsOfFollows = append(userIdsOfFollows, follow.ID)
		}

		/** TODO: This is kind of sub-optimal, as we fetch both 25 hobbits as well as records. Optimaly, we would only fetch
		the 25 items we need from the db. This would also allow us to continue the feed with pages.
		*/
		// First, fetch all records of people you follow
		recentRecordsOfFollowers := []*models.NumericRecord{}
		err = db.Preload("Hobbit.User").Joins("Hobbit").Joins("LEFT JOIN Users on hobbit.user_id = users.id").Where("hobbit.user_id IN ?", userIdsOfFollows).Limit(25).Order("numeric_records.created_at DESC").Find(&recentRecordsOfFollowers).Error
		if err != nil {
			slog.Error("failed to find records of people you follow", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Then we fetch all of the hobbits of people you follow
		recentHobbitsOfFollowers := []*models.Hobbit{}
		err = db.Joins("User").Where("user_id in ?", userIdsOfFollows).Limit(25).Order("created_at DESC").Find(&recentHobbitsOfFollowers).Error

		relevantEvents := []FeedEvent{}

		for _, r := range recentRecordsOfFollowers {
			relevantEvents = append(relevantEvents, FeedEvent{
				FeedEventTypus: FeedEventTypusRecordCreated,
				CreatedAt:      r.CreatedAt,
				Payload:        r,
			})
		}
		for _, h := range recentHobbitsOfFollowers {
			relevantEvents = append(relevantEvents, FeedEvent{
				FeedEventTypus: FeedEventTypusHobbitCreated,
				CreatedAt:      h.CreatedAt,
				Payload:        h,
			})
		}

		sort.Slice(relevantEvents, func(i, j int) bool {
			return relevantEvents[i].CreatedAt.After(relevantEvents[j].CreatedAt)
		})

		upperMax := min(len(relevantEvents)-1, 25)

		relevantEvents = relevantEvents[0:upperMax]

		err = json.NewEncoder(w).Encode(relevantEvents)
		if err != nil {
			slog.Error("failed to encode response", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
