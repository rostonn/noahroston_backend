package oauth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rostonn/noahroston_backend/app"
	"github.com/rostonn/noahroston_backend/models"
	"go.uber.org/zap"
)

func LoginWithAmazon(code string, config app.Configuration) (*models.User, *OauthError) {
	zap.S().Debug("loginWithAmazon: amazon token", zap.String("code", code))
	user := &models.User{}
	oauthError := &OauthError{}

	amazonTokenRequestMap := map[string]interface{}{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     config.AmazonClientId,
		"client_secret": config.AmazonClientSecret,
		"redirect_uri":  config.AmazonRedirectURI,
	}

	jsonTokenRequestBytes, _ := json.Marshal(amazonTokenRequestMap)

	req, err := http.NewRequest("POST", config.AmazonAccessTokenURL, bytes.NewBuffer(jsonTokenRequestBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error("Post to amazon access token url")
		// panic(err)
		oauthError.Code = 500
		oauthError.Message = "Post to amazon access token url"
		oauthError.Error = err
		return nil, oauthError
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var amazonResponseMap map[string]string
	json.Unmarshal(body, &amazonResponseMap)

	accessToken, ok := amazonResponseMap["access_token"]
	if !ok || accessToken == "" {
		oauthError.Message = "No access token from Amazon"
		oauthError.Code = 401
		return nil, oauthError
	}

	zap.S().Info("access", accessToken)
	zap.S().Info(amazonResponseMap)

	// Make request to get user info
	userInfoReq, err := http.NewRequest("GET", config.AmazonProfileURL, nil)
	userInfoReq.Header.Set("Accept", "application/json")
	userInfoReq.Header.Set("x-amz-access-token", accessToken)

	zap.S().Debug("Amazon Acess Token: " + accessToken)
	if err != nil {
		zap.S().Error("Amazon access error", zap.Error(err))
		oauthError.Error = err
		oauthError.Code = 500
		oauthError.Message = "Amazon unmarshall error"
		return nil, oauthError
	}

	resp, err = client.Do(userInfoReq)
	if err != nil {
		zap.S().Error("Amazon Request Failed", zap.Error(err))
		oauthError.Error = err
		oauthError.Code = 500
		oauthError.Message = "Amazon unmarshall error"
		return nil, oauthError
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	var userObj map[string]interface{}
	err = json.Unmarshal(body, &userObj)

	b, err := json.MarshalIndent(userObj, "", "  ")
	zap.S().Info("AMZN User " + string(b))
	zap.S().Info("AMZN User", string(b))

	user.LastOauth = "AMAZON"

	email, _ := userObj["email"].(string)
	user.Email = email

	return user, nil
}
