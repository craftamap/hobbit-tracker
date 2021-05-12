package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/craftamap/hobbit-tracker/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var diskMode bool
var port int

//go:embed frontend/dist
var content embed.FS
var db *gorm.DB

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	flag.BoolVar(&diskMode, "disk-mode", false, "disk mode")
	flag.IntVar(&port, "port", 8080, "port")
	flag.Parse()
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		next.ServeHTTP(w, r)
	})
}

func frontendHandler() (http.Handler, error) {
	var fsys fs.FS
	fsys = fs.FS(content)
	contentStatic, err := fs.Sub(fsys, "frontend/dist")
	if err != nil {
		return nil, err
	}
	if diskMode {
		log.Warn("Disk Mode")
		contentStatic = os.DirFS("frontend/dist")
	}
	return http.FileServer(http.FS(contentStatic)), nil
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("hobbits.sqlite"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Info("AutoMigrating DB")
	db.AutoMigrate(&models.Hobbit{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.NumericRecord{})
	log.Info("AutoMigrated DB")

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Use(loggingMiddleware)
	r.Use(AuthToContextMiddleBuilder(db, log))

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", BuildHandleLogin(db, log)).Methods("POST")
	auth.HandleFunc("/logout", BuildHandleLogout(log))

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth", BuildHandleAPIGetAuth(db, log)).Methods("GET")

	hobbits := api.PathPrefix("/hobbits").Subrouter()
	hobbits.Handle("/", AuthMiddlewareBuilder(log)(
		http.HandlerFunc(BuildHandleAPIPostHobbit(db, log)),
	)).Methods("POST")
	hobbits.Handle("/{id:[0-9]+}", AuthMiddlewareBuilder(log)(
		http.HandlerFunc(BuildHandleAPIPutHobbit(db, log)),
	)).Methods("PUT")
	hobbits.Handle("/{id:[0-9]+}", BuildHandleAPIGetHobbit(db, log)).Methods("GET")
	hobbits.Handle("/", BuildHandleAPIGetHobbits(db, log)).Methods("GET")

	records := hobbits.PathPrefix("/{hobbit_id:[0-9]+}/records").Subrouter()
	records.Handle("/", BuildHandleAPIGetRecords(db, log)).Methods("GET")
	records.Handle("/", AuthMiddlewareBuilder(log)(
		http.HandlerFunc(BuildHandleAPIPostRecord(db, log)),
	)).Methods("POST")
	records.Handle("/{record_id:[0-9]+}", AuthMiddlewareBuilder(log)(
		http.HandlerFunc(BuildHandleAPIPutRecord(db, log)),
	)).Methods("PUT")
	records.Handle("/{record_id:[0-9]+}", AuthMiddlewareBuilder(log)(
		http.HandlerFunc(BuildHandleAPIDeleteRecord(db, log)),
	)).Methods("DELETE")
	records.Handle("/heatmap", BuildHandleAPIGetRecordsForHeatmap(db, log)).Methods("GET")

	profile := api.PathPrefix("/profile").Subrouter()
	profileMe := profile.PathPrefix("/me").Subrouter()
	profileMe.Handle("/", BuildHandleAPIGetAuth(db, log))
	profileMe.Handle("/hobbits", AuthMiddlewareBuilder(log)(
		http.HandlerFunc(BuildHandleAPIProfileGetHobbits(db, log)),
	)).Methods("GET")

	frontend, err := frontendHandler()
	if err != nil {
		fmt.Println(err)
		return
	}
	r.PathPrefix("/").Handler(frontend)
	listeningOn := fmt.Sprintf(":%d", port)
	log.Infof("Listening on %s", listeningOn)
	http.ListenAndServe(listeningOn, r)
}
