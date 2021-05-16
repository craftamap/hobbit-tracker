package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	routesApiHobbits "github.com/craftamap/hobbit-tracker/routes/api/hobbits"
	routesApiProfile "github.com/craftamap/hobbit-tracker/routes/api/profile"
	routesApiRecords "github.com/craftamap/hobbit-tracker/routes/api/records"

	middlewareAuth "github.com/craftamap/hobbit-tracker/middleware/auth"
	"github.com/craftamap/hobbit-tracker/middleware/authToContext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/craftamap/hobbit-tracker/routes"
	routesApiAuth "github.com/craftamap/hobbit-tracker/routes/api/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var diskMode bool
var port int

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("jMcBBEBKAzw89XNb")
	Store = sessions.NewCookieStore(key)
)

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
	var fsys = fs.FS(content)
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
	db.AutoMigrate(&models.AppPassword{})
	log.Info("AutoMigrated DB")

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Use(loggingMiddleware)
	r.Use(authToContext.New(db, log, Store))

	authMiddlewareBuilder := middlewareAuth.Builder(log)

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", routes.BuildHandleLogin(db, log, Store)).Methods("POST")
	auth.HandleFunc("/logout", routes.BuildHandleLogout(log, Store))

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth", routesApiAuth.BuildHandleAPIGetAuth(db, log)).Methods("GET")

	hobbits := api.PathPrefix("/hobbits").Subrouter()
	hobbits.Handle("/", authMiddlewareBuilder.Build(
		http.HandlerFunc(routesApiHobbits.BuildHandleAPIPostHobbit(db, log)),
	)).Methods("POST")
	hobbits.Handle("/{id:[0-9]+}", authMiddlewareBuilder.Build(
		http.HandlerFunc(routesApiHobbits.BuildHandleAPIPutHobbit(db, log)),
	)).Methods("PUT")
	hobbits.Handle("/{id:[0-9]+}", routesApiHobbits.BuildHandleAPIGetHobbit(db, log)).Methods("GET")
	hobbits.Handle("/", routesApiHobbits.BuildHandleAPIGetHobbits(db, log)).Methods("GET")

	records := hobbits.PathPrefix("/{hobbit_id:[0-9]+}/records").Subrouter()
	records.Handle("/", routesApiRecords.BuildHandleAPIGetRecords(db, log)).Methods("GET")
	records.Handle("/", authMiddlewareBuilder.Build(
		http.HandlerFunc(routesApiRecords.BuildHandleAPIPostRecord(db, log)),
	)).Methods("POST")
	records.Handle("/{record_id:[0-9]+}", authMiddlewareBuilder.Build(
		http.HandlerFunc(routesApiRecords.BuildHandleAPIPutRecord(db, log)),
	)).Methods("PUT")
	records.Handle("/{record_id:[0-9]+}", authMiddlewareBuilder.Build(
		http.HandlerFunc(routesApiRecords.BuildHandleAPIDeleteRecord(db, log)),
	)).Methods("DELETE")
	records.Handle("/heatmap", routesApiRecords.BuildHandleAPIGetRecordsForHeatmap(db, log)).Methods("GET")

	profile := api.PathPrefix("/profile").Subrouter()
	profileMe := profile.PathPrefix("/me").Subrouter()
	profileMe.Use(authMiddlewareBuilder.Build)
	profileMe.Handle("/", routesApiAuth.BuildHandleAPIGetAuth(db, log))
	profileMe.Handle("/hobbits", http.HandlerFunc(routesApiProfile.BuildHandleAPIProfileGetHobbits(db, log))).Methods("GET")
	profileMeAppPassword := profileMe.PathPrefix("/apppassword").Subrouter()
	profileMeAppPassword.Use(authMiddlewareBuilder.WithPermitAppPasswordAuth(false).Build)
	profileMeAppPassword.HandleFunc("/", routesApiProfile.BuildHandleAPIProfileGetAppPasswords(db, log)).Methods("GET")
	profileMeAppPassword.HandleFunc("/", routesApiProfile.BuildHandleAPIProfilePostAppPassword(db, log)).Methods("POST")

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
