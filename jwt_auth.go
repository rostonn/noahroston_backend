package main

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/rostonn/noahroston_backend/models"
)

type CustomClaims struct {
	User            models.User            `json:"user"`
	userLoginRecord models.UserLoginRecord `json:"userLoginRecord"`
	jwt.StandardClaims
}

func (a *App) createAndReturnJWT(w http.ResponseWriter, r *http.Request, user models.User) {

	ipAddress := getIPAdress(r)

	fmt.Println("IP", ipAddress)

	userLoginRecord := getIpAddressInfo(a.config.IpStackApiKey, ipAddress)
	userLoginRecord.OauthProvider = user.LastOauth
	userLoginRecord.UserID = user.ID

	err := userLoginRecord.CreateUserLoginRecord(a.DB)
	if err != nil {
		panic(err)
	}

	fmt.Println("Creating JWT ...")

	claims := CustomClaims{
		user,
		userLoginRecord,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "noahroston",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err := token.SignedString(a.config.PrivateKey)
	fmt.Printf("%v %v", ss, err)

	if err != nil {
		fmt.Println("JWT ERROR")
		panic(err)
		respondWithError(w, 500, "Error Creating JWT")
	}
	w.WriteHeader(200)
	w.Write([]byte(ss))
}