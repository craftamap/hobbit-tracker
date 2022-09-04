package main

import (
	"context"
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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	flagDiskMode bool
	flagPort     int
	flagVerbose  bool
)

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

	flag.BoolVar(&flagDiskMode, "disk-mode", false, "disk mode")
	flag.IntVar(&flagPort, "port", 8080, "port")
	flag.BoolVar(&flagVerbose, "v", false, "verbose, enables debug logs")
	flag.Parse()

	if flagVerbose {
		log.SetLevel(logrus.DebugLevel)
	}
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
	if flagDiskMode {
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

func newExporter() (trace.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://172.18.0.2:14268/api/traces")))
}
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fib"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

func main() {
	var err error

	exp, err := newExporter()
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	gormConfig := &gorm.Config{}
	if flagVerbose {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	}

	db, err = gorm.Open(sqlite.Open("hobbits.sqlite"), gormConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Info("Migrating DB")
	// Manual Migration
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
	log.Info("Checking for manual migrations")
	{
		// Migration from v0.2.1 to above
		records := []models.NumericRecord{}
		db.Where("created_at IS NULL").Find(&records)
		if len(records) > 0 {
			log.Info("Found migration from v0.2.1, performing migration")
			db.Model(&models.NumericRecord{}).Where("created_at IS NULL").Updates(&models.NumericRecord{CreatedAt: time.Now()})

			db.Model(&models.Hobbit{}).Where("created_at IS NULL").Updates(&models.Hobbit{CreatedAt: time.Now()})
			log.Info("Found migration from v0.2.1, done")
		}
	}
	log.Info("Migrated DB")

	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key := []byte("jMcBBEBKAzw89XNb")
	Store = gormstore.New(db, key)
	// db cleanup every hour
	// close quit channel to stop cleanup
	quit := make(chan struct{})
	go Store.PeriodicCleanup(1*time.Hour, quit)
	defer close(quit)
	eventHub := hub.New(log)
	eventHub.Run()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newCtx, span := otel.Tracer("hobbit").Start(r.Context(), "Request")
			span.SetAttributes(attribute.String("Method", r.Method))
			span.SetAttributes(attribute.String("Path", r.URL.Path))
			next.ServeHTTP(w, r.WithContext(newCtx))
			span.End()
		})
	})
	r.Use(loggingMiddleware)

	r.Use(handlers.RecoveryHandler(handlers.RecoveryLogger(&customRecoveryLogger{log})))
	r.Use(requestcontext.New(Store, db, log, eventHub))
	r.Use(authtocontext.New())

	routes.RegisterRoutes(r)
	websockets.RegisterRoutes(r)

	frontend, err := frontendHandler()
	if err != nil {
		fmt.Println(err)
		return
	}
	r.PathPrefix("/").Handler(frontend)
	listeningOn := fmt.Sprintf(":%d", flagPort)
	log.Infof("Listening on %s", listeningOn)
	err = http.ListenAndServe(listeningOn, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
