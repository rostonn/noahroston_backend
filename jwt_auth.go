package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/rostonn/noahroston_backend/models"
)

type CustomClaims struct {
	User *models.User `json:"user"`
	jwt.StandardClaims
}

func (a *App) validateJWT(w http.ResponseWriter, r *http.Request, token string) {
	zap.S().Debug("Validating JWT " + token)

	var myClaims CustomClaims
	t, err := jwt.ParseWithClaims(token, &myClaims, func(token *jwt.Token) (interface{}, error) {
		return a.config.PublicKey, nil
	})

	if err != nil {
		zap.S().Error(err)
		zap.S().Error("Token Error")
		respondWithError(w, 401, "Error Creating JWT")
		return
	}

	zap.S().Debug(t.Claims)

	w.Write(nil)
}

func (a *App) createAndReturnJWT(w http.ResponseWriter, r *http.Request, user *models.User) {

	ipAddress := getIPAdress(r)

	zap.S().Info("IP " + ipAddress)

	userLoginRecord := getIpAddressInfo(a.config.IpStackApiKey, ipAddress)

	userLoginRecord.IpAddress = ipAddress

	userLoginRecord.OauthProvider = user.LastOauth
	userLoginRecord.UserID = user.ID

	err := userLoginRecord.CreateUserLoginRecord(a.DB)
	if err != nil {
		panic(err)
	}

	zap.S().Debug("Creating JWT")
	user.UserLoginRecord = userLoginRecord

	claims := CustomClaims{
		user,

		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 3).Unix(),
			Issuer:    "noahroston",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err := token.SignedString(a.config.PrivateKey)
	ssString := fmt.Sprintf("%v %v", ss, err)
	zap.S().Debug(ssString)

	if err != nil {
		zap.S().Error("JWT ERROR")
		respondWithError(w, 500, "Error Creating JWT")
	}
	w.WriteHeader(200)
	w.Write([]byte(ss))
}
