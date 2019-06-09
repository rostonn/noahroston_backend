package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rostonn/noahroston_backend/app"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	config app.Configuration
	Zlog   *zap.Logger
}

func (a *App) Initialize(config app.Configuration, logger *zap.Logger) {
	a.Zlog = logger
	connectionString := fmt.Sprintf("%s:%s@/%s", config.DbUsername, config.DbPassword, config.Dbname)
	var err error

	a.Zlog.Debug("DB INFO",
		zap.String("Username", config.DbUsername),
		zap.String("DBName", config.Dbname))

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
