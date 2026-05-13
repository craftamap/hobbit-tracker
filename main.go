package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/craftamap/hobbit-tracker/hub"
	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/craftamap/hobbit-tracker/routes"
	"github.com/craftamap/hobbit-tracker/routes/api"
	"github.com/craftamap/hobbit-tracker/websockets"
	"github.com/gorilla/handlers"
	"github.com/wader/gormstore/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	flagDiskMode bool
	flagPort     int
	flagVerbose  bool
	flagDatabase string
)

var (
	// Store represents the Cookie Store
	Store *gormstore.Store
)

//go:embed frontend/dist
var content embed.FS
var db *gorm.DB

func init() {
	flag.BoolVar(&flagDiskMode, "disk-mode", false, "disk mode")
	flag.IntVar(&flagPort, "port", 8080, "port")
	flag.BoolVar(&flagVerbose, "v", false, "verbose, enables debug logs")
	flag.StringVar(&flagDatabase, "db", "hobbits.sqlite", "path to the database")
	flag.Parse()

	if flagVerbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "addr", r.RemoteAddr, "method", r.Method, "url", r.URL)
		next.ServeHTTP(w, r)
	})
}

func frontendHandler() (http.Handler, error) {
	var fsys = fs.FS(content)
	contentStatic, err := fs.Sub(fsys, "frontend/dist")
	if err != nil {
		return nil, err
	}
	if flagDiskMode {
		slog.Warn("Disk Mode")
		contentStatic = os.DirFS("frontend/dist")
	}
	return http.FileServer(http.FS(contentStatic)), nil
}

type customRecoveryLogger struct {
}

func (c *customRecoveryLogger) Println(msgs ...any) {
	slog.Error("recovering", "msgs", msgs)
}

func main() {
	var err error

	gormConfig := &gorm.Config{}
	if flagVerbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	db, err = gorm.Open(sqlite.Open(flagDatabase), gormConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	slog.Info("Migrating DB")
	// Manual Migration
	slog.Info("AutoMigrating DB")
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

	slog.Info("AutoMigrated DB")
	slog.Info("Checking for manual migrations")
	{
		// Migration from v0.2.1 to above
		records := []models.NumericRecord{}
		db.Where("created_at IS NULL").Find(&records)
		if len(records) > 0 {
			slog.Info("Found migration from v0.2.1, performing migration")
			db.Model(&models.NumericRecord{}).Where("created_at IS NULL").Updates(&models.NumericRecord{CreatedAt: time.Now()})

			db.Model(&models.Hobbit{}).Where("created_at IS NULL").Updates(&models.Hobbit{CreatedAt: time.Now()})
			slog.Info("Found migration from v0.2.1, done")
		}
	}
	slog.Info("Migrated DB")

	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key := []byte("jMcBBEBKAzw89XNb")
	Store = gormstore.New(db, key)
	// db cleanup every hour
	// close quit channel to stop cleanup
	quit := make(chan struct{})
	go Store.PeriodicCleanup(1*time.Hour, quit)
	defer close(quit)
	eventHub := hub.New()
	eventHub.Run()

	rootRouter := http.NewServeMux()

	rootRouter.Handle("POST /auth/login", routes.BuildHandleLogin())
	rootRouter.Handle("POST /auth/logout", routes.BuildHandleLogout())
	rootRouter.Handle("/api/", api.GetRoutes())
	// routes.GetRoutes(db, log, Store)
	rootRouter.Handle("/ws", websockets.GetRoutes())

	frontend, err := frontendHandler()
	if err != nil {
		fmt.Println(err)
		return
	}
	rootRouter.Handle("/", frontend)

	r := authtocontext.New()(rootRouter)
	r = requestcontext.New(Store, db, eventHub)(r)
	r = handlers.RecoveryHandler(handlers.RecoveryLogger(&customRecoveryLogger{}))(r)
	r = loggingMiddleware(r)

	listeningOn := fmt.Sprintf(":%d", flagPort)
	slog.Info("Listening on", "port", flagPort)
	err = http.ListenAndServe(listeningOn, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
