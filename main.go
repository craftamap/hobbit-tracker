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

var DiskMode bool

//go:embed frontend/dist
var content embed.FS
var db *gorm.DB

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	flag.BoolVar(&DiskMode, "disk-mode", false, "disk mode")
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
	if DiskMode {
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

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", BuildHandleLogin(db, log)).Methods("POST")
	auth.HandleFunc("/logout", BuildHandleLogout(log))

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth", BuildHandleApiGetAuth(db, log)).Methods("GET")
	hobbits := api.PathPrefix("/hobbits").Subrouter()

	hobbits.Handle("/", AuthMiddlewareBuilder(log)(
		http.HandlerFunc(BuildHandleApiPostHobbit(db, log)),
	)).Methods("POST")

	hobbits.Handle("/{id:[0-9]+}", AuthMiddlewareBuilder(log)(
		http.HandlerFunc(BuildHandleApiPutHobbit(db, log)),
	)).Methods("PUT")

	hobbits.Handle("/{id:[0-9]+}", BuildHandleApiGetHobbit(db, log)).Methods("GET")
	hobbits.Handle("/", BuildHandleApiGetHobbits(db, log)).Methods("GET")

	records := hobbits.PathPrefix("/{hobbit_id:[0-9]+}/records").Subrouter()
	records.Handle("/", BuildHandleApiGetRecords(db, log)).Methods("GET")
	records.Handle("/", AuthMiddlewareBuilder(log)(
		http.HandlerFunc(BuildHandleApiPostRecord(db, log)),
	)).Methods("POST")

	frontend, err := frontendHandler()
	if err != nil {
		fmt.Println(err)
		return
	}
	r.PathPrefix("/").Handler(frontend)
	log.Info("Listening on :3000")
	http.ListenAndServe(":3000", r)
}
