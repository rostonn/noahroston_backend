package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rostonn/noahroston_backend/models"

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
		return
	}
	var m map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		respondWithError(w, 400, "Bad Request - Body must be json")
		return
	}

	code := m["code"]
	if code == "" {
		respondWithError(w, 400, "Bad Request - Code must not be empty")
		return
	}
	zap.S().Info("LoginUser Code", zap.String("code", code))

	switch provider {
	case "amazon":
		a.loginWithAmazon(w, r, code)
	case "google":
		a.loginWithGoogle(w, r, code)
	case "facebook":
		a.loginWithFacebook(w, r, code)
	case "test":
		a.loginWithTester(w, r)
	default:
		respondWithError(w, 400, "Bad Request - Provider "+provider+" unknown")
	}
}

func (a *App) loginWithTester(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	user.LastOauth = "TEST"
	user.Email = "tester@test.com"
	a.logUserIn(w, r, user, nil)
}

func (a *App) loginWithAmazon(w http.ResponseWriter, r *http.Request, code string) {
	user, userError := oauth.LoginWithAmazon(code, a.config)
	a.logUserIn(w, r, user, userError)
}

func (a *App) loginWithGoogle(w http.ResponseWriter, r *http.Request, code string) {
	user, userError := oauth.LoginWithGoogle(code, a.config)
	a.logUserIn(w, r, user, userError)
}

func (a *App) logUserIn(w http.ResponseWriter, r *http.Request, user *models.User, userError *oauth.OauthError) {
	if userError != nil {
		zap.S().Error("Login Error")
		respondWithError(w, userError.Code, userError.Message)
	}

	err := user.LoginUser(a.DB)
	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	zap.S().Debug("UserID: " + string(user.ID))
	a.createAndReturnJWT(w, r, user)
}

func (a *App) loginWithFacebook(w http.ResponseWriter, r *http.Request, code string) {
	user, userError := oauth.LoginWithFacebook(code, a.config)
	a.logUserIn(w, r, user, userError)
}

func (a *App) checkToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

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

	zap.S().Debug("Token " + token)

	a.validateJWT(w, r, token)

}
