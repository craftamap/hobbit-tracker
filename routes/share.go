package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/gorilla/mux"
)

func HandleShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := requestcontext.Log(r)
		db := requestcontext.DB(r)

		user := models.User{}

		err := db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		byteData, _ := ioutil.ReadAll(r.Body)
		log.Info(byteData)

		r.ParseMultipartForm(5 * 1024 * 1024 * 1024)
		log.WithField("files", r.MultipartForm).Info()
		file, header, err := r.FormFile("file")
		if err != nil {
			log.WithField("error", err).Error()
			w.WriteHeader(500)
		}
		// log.WithFields(logrus.Fields{"file": file, "header": header}).Info()
		contentType := header.Header["Content-Type"][0]
		if contentType != "application/gpx+xml" {
			log.WithField("contentType", contentType).Warn("contentType is not matching")
		}

		buf := &bytes.Buffer{}
		buf.ReadFrom(file)
		file.Close()

		// TODO: also store a user-facing uuid, as if all records are deleted gorm/sqlite restart with id 1
		shareFile := &models.TemporaryShareFile{
			Bytes:    buf.Bytes(),
			MimeType: contentType,
			User:     user,
		}
		db.Create(shareFile)
		log.WithField("id", shareFile.ID).Info()

		http.Redirect(w, r, fmt.Sprintf("/#/share/%d", shareFile.ID), http.StatusTemporaryRedirect)
	}
}

func HandleGetShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := requestcontext.Log(r)
		db := requestcontext.DB(r)

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
		numericID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shareFile := &models.TemporaryShareFile{}
		err = db.Where(&models.TemporaryShareFile{ID: uint(numericID), UserID: user.ID}).Joins("User").First(shareFile).Error
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		gpxStats, err := ParseGpx(shareFile.Bytes)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(gpxStats)
	}
}

type GpxStats struct {
	Time       time.Time
	MovingTime time.Duration
	Distance   int
}

func ParseGpx(byteData []byte) (interface{}, error) {
	parser, err := xmlquery.Parse(bytes.NewReader(byteData))
	if err != nil {
		return nil, err
	}

	extensionData, err := xmlquery.Query(parser, "//gpxtrkx:TrackStatsExtension")
	if err != nil {
		return nil, err
	}

	movingTimeNode, err := xmlquery.Query(extensionData, "//gpxtrkx:MovingTime")
	if err != nil {
		return nil, err
	}
	movingTime, err := strconv.ParseFloat(movingTimeNode.InnerText(), 64)
	if err != nil {
		return nil, err
	}

	distanceNode, err := xmlquery.Query(extensionData, "//gpxtrkx:Distance")
	if err != nil {
		return nil, err
	}
	distance, err := strconv.ParseFloat(distanceNode.InnerText(), 64)
	if err != nil {
		return nil, err
	}

	firstTimestampNode, err := xmlquery.Query(parser, "//trkseg//trkpt//time")
	if err != nil {
		return nil, err
	}
	firstTimestamp, err := time.Parse(time.RFC3339, firstTimestampNode.InnerText())
	if err != nil {
		return nil, err
	}

	return GpxStats{
		Distance:   int(distance),
		MovingTime: time.Duration(movingTime * 1000 * 1000 * 1000),
		Time:       firstTimestamp,
	}, nil
}
