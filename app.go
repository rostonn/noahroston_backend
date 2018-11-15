package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang/glog"
	log "github.com/sirupsen/logrus"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	config Configuration
}

func (a *App) Initialize(config Configuration) {
	connectionString := fmt.Sprintf("%s:%s@/%s", config.DbUsername, config.DbPassword, config.Dbname)
	var err error

	glog.Info("Prepare to repel boarders")

	// sugar := zap.NewDevelopment()
	// defer sugar.Sync()
	// sugar.Infow("failed to fetch URL",
	// 	"url", "http://example.com",
	// 	"attempt", 3,
	// 	"backoff", time.Second,
	// )
	// sugar.Infof("failed to fetch URL: %s", "http://example.com")

	logger := &log.Logger{
		Out:   os.Stderr,
		Level: log.DebugLevel,
		Formatter: &log.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
		ReportCaller: true,
	}
	// logger.SetReportCaller(true)

	logger.WithFields(log.Fields{
		"db user":  config.DbUsername,
		"password": config.DbPassword,
		"DB":       config.Dbname,
	}).Error("Database Connection")

	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()

	a.initializeRoutes()
	a.config = config
}

func (a *App) Run(addr string) {
	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(addr, handlers.CORS(headersOk, originsOk, methodsOk)(a.Router)))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
