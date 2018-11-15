package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rostonn/noahroston_backend/models"
)

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
	fmt.Println("Code", code)
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
	fmt.Println("loginWithAmazon: amazon token", code)

	amazonTokenRequestMap := map[string]interface{}{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     a.config.AmazonClientId,
		"client_secret": a.config.AmazonClientSecret,
		"redirect_uri":  a.config.AmazonRedirectURI,
	}

	jsonTokenRequestBytes, _ := json.Marshal(amazonTokenRequestMap)

	req, err := http.NewRequest("POST", a.config.AmazonAccessTokenURL, bytes.NewBuffer(jsonTokenRequestBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var amazonResponseMap map[string]string
	json.Unmarshal(body, &amazonResponseMap)

	accessToken, ok := amazonResponseMap["access_token"]
	if !ok {
		respondWithError(w, 401, "No access token from Amazon")
	}
	// Make request to get user info
	userInfoReq, err := http.NewRequest("GET", a.config.AmazonProfileURL, nil)
	userInfoReq.Header.Set("Accept", "application/json")
	userInfoReq.Header.Set("x-amz-access-token", accessToken)

	fmt.Println("Access Token:", accessToken)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, 500, "Amazon unmarshall error")
	}

	resp, err = client.Do(userInfoReq)
	if err != nil {
		fmt.Println("Amazon Request Failed")
		// panic(err)
		respondWithError(w, 500, "Amazon unmarshall error")
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	var userObj map[string]interface{}
	err = json.Unmarshal(body, &userObj)

	b, err := json.MarshalIndent(userObj, "", "  ")
	fmt.Println("AMZN User", string(b))

	var user models.User
	user.LastOauth = "AMAZON"

	email, _ := userObj["email"].(string)

	user.Email = email

	err = user.LoginUser(a.DB)
	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	fmt.Println(user)
	// user object is filled out

	a.createAndReturnJWT(w, r, user)
}
