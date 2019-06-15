package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
		a.Zlog.Error("Returning error from loginWithAmazon")
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

func (a *App) checkToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Checking token request type: " + r.Method)

	authHeader := r.Header.Get("Authorization")

	fmt.Println("Authorization Header", authHeader)

	if authHeader == "" {
		respondWithError(w, 401, "UNAUTHORIZED")
	}

	authArr := strings.Split(authHeader, " ")

	if len(authArr) != 2 {
		respondWithError(w, 400, "BAD REQUEST")
	}

	if authArr[0] != "Bearer" {
		respondWithError(w, 400, "BAD REQUEST")
	}

	token := authArr[1]

	a.Zlog.Debug("Token " + token)

	a.validateJWT(w, r, token)

}
