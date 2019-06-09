package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rostonn/noahroston_backend/oauth"
	"go.uber.org/zap"
)

func (a *App) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (a *App) loginUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	if provider == "" {
		respondWithError(w, 400, "Bad Request - Provider cannot be nil")
	}
	var m map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		respondWithError(w, 400, "Bad Request - Body must be json")
	}

	code := m["code"]
	if code == "" {
		respondWithError(w, 400, "Bad Request - Code must not be empty")
	}
	a.Zlog.Info("LoginUser Code", zap.String("code", code))

	switch provider {
	case "amazon":
		a.loginWithAmazon(w, r, code)
	default:
		fmt.Println("Goes Here respond with error?")
		respondWithError(w, 400, "Bad Request - Provider "+provider+" unknown")
	}
	// Switch provider and forward request and response on

}

func (a *App) loginWithAmazon(w http.ResponseWriter, r *http.Request, code string) {
	user, userError := oauth.LoginWithAmazon(code, a.config, a.Zlog)
	if userError != nil {
		respondWithError(w, userError.Code, userError.Message)
	}

	err := user.LoginUser(a.DB)
	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	a.Zlog.Debug("UserID: " + string(user.ID))
	fmt.Println(user)

	a.createAndReturnJWT(w, r, user)
}
