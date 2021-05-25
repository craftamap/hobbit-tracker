package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/craftamap/hobbit-tracker/hub"
	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/craftamap/hobbit-tracker/routes"
	"github.com/craftamap/hobbit-tracker/websockets"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/wader/gormstore/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var diskMode bool
var port int

var (
	// Store represents the Cookie Store
	Store *gormstore.Store
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

type customRecoveryLogger struct {
	log *logrus.Logger
}

func (c *customRecoveryLogger) Println(msgs ...interface{}) {
	c.log.Errorln(msgs...)
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("hobbits.sqlite"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Info("AutoMigrating DB")
	err = db.AutoMigrate(&models.Hobbit{})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.AutoMigrate(&models.NumericRecord{})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.AutoMigrate(&models.AppPassword{})
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Info("AutoMigrated DB")

	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key := []byte("jMcBBEBKAzw89XNb")
	Store = gormstore.New(db, key)
	// db cleanup every hour
	// close quit channel to stop cleanup
	quit := make(chan struct{})
	go Store.PeriodicCleanup(1*time.Hour, quit)
	defer close(quit)
	eventHub := hub.New()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Use(loggingMiddleware)

	// TODO: find a way to make
	r.Use(handlers.RecoveryHandler(handlers.RecoveryLogger(&customRecoveryLogger{log})))
	r.Use(requestcontext.New(Store, db, log, eventHub))
	r.Use(authtocontext.New())

	routes.RegisterRoutes(r, db, log, Store)
	websockets.RegisterRoutes(r, db, log, Store, eventHub)

	tckr := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-tckr.C:
				eventHub.Broadcast(hub.ServerSideEvent{
					Typus: hub.HobbitCreated,
				})
			}
		}
	}()

	frontend, err := frontendHandler()
	if err != nil {
		fmt.Println(err)
		return
	}
	r.PathPrefix("/").Handler(frontend)
	listeningOn := fmt.Sprintf(":%d", port)
	log.Infof("Listening on %s", listeningOn)
	err = http.ListenAndServe(listeningOn, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
