package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/craftamap/hobbit-tracker/middleware/authToContext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/craftamap/hobbit-tracker/routes"
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

	routes.RegisterRoutes(r, db, log, Store)

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
