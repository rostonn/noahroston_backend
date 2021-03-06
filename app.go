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
}

func (a *App) Initialize(config app.Configuration) {
	// db, err := sql.Open("mysql", "db_user:password@tcp(localhost:3306)/my_db")
	connectionString := fmt.Sprintf("%s:%s@%s", config.DbUsername, config.DbPassword, config.Dbname)
	var err error

	zap.S().Debug("DB INFO ",config.DbUsername," ",config.Dbname)

	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = a.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
	a.config = config
}

func (a *App) Run(addr string) {
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
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
