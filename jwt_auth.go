package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/rostonn/noahroston_backend/models"
)

type CustomClaims struct {
	User            *models.User           `json:"user"`
	UserLoginRecord models.UserLoginRecord `json:"userLoginRecord"`
	jwt.StandardClaims
}

func (a *App) validateJWT(w http.ResponseWriter, r *http.Request, token string) {
	a.Zlog.Debug("Validating JWT " + token)

	// fmt.Println(a.config.PrivateKey)

	// fmt.Println("Public Key Bytes")
	// fmt.Println(a.config.PublicKeyBytes)

	var myClaims CustomClaims
	t, err := jwt.ParseWithClaims(token, &myClaims, func(token *jwt.Token) (interface{}, error) {
		return a.config.PublicKey, nil
	})

	// fmt.Println("claims ...", t)

	if err != nil {
		fmt.Println(err)
		a.Zlog.Error("Token Error")
		respondWithError(w, 401, "Error Creating JWT")
		return
	}
	// validErr := t.Claims.Valid
	// if validErr != nil {
	// 	fmt.Println("Invalid JWT??")
	// 	fmt.Println(validErr)
	// 	respondWithError(w, 401, "Invalid JWT")
	// 	return
	// }

	fmt.Println("Expires At", myClaims.ExpiresAt)
	fmt.Println("Check Expires", myClaims.VerifyExpiresAt)
	fmt.Println(myClaims.User.Email)
	fmt.Println(myClaims.User.LastOauth)
	fmt.Println("Token is valid responding")

	fmt.Println(t.Claims)

	// w.WriteHeader(200)
	w.Write(nil)
}

func (a *App) createAndReturnJWT(w http.ResponseWriter, r *http.Request, user *models.User) {

	ipAddress := getIPAdress(r)

	a.Zlog.Info("IP " + ipAddress)

	userLoginRecord := getIpAddressInfo(a.config.IpStackApiKey, ipAddress)
	userLoginRecord.OauthProvider = user.LastOauth
	userLoginRecord.UserID = user.ID

	err := userLoginRecord.CreateUserLoginRecord(a.DB)
	if err != nil {
		panic(err)
	}

	a.Zlog.Debug("Creating JWT")

	claims := CustomClaims{
		user,
		userLoginRecord,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 3).Unix(),
			Issuer:    "noahroston",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err := token.SignedString(a.config.PrivateKey)
	ssString := fmt.Sprintf("%v %v", ss, err)
	a.Zlog.Debug(ssString)

	if err != nil {
		a.Zlog.Error("JWT ERROR")
		respondWithError(w, 500, "Error Creating JWT")
	}
	w.WriteHeader(200)
	w.Write([]byte(ss))
}
