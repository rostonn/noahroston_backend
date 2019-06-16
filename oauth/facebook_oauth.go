package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rostonn/noahroston_backend/app"
	"github.com/rostonn/noahroston_backend/models"
	"go.uber.org/zap"
)

func LoginWithFacebook(code string, config app.Configuration, logger *zap.Logger) (*models.User, *OauthError) {
	logger.Debug("loginWithFacebook: token", zap.String("code", code))
	user := &models.User{}
	oauthError := &OauthError{}

	resp, err := http.Get("https://graph.facebook.com/v3.3/me?fields=email&access_token=" + code)

	if err != nil {
		logger.Error("Get to facebook token url")
		// panic(err)
		oauthError.Code = 500
		oauthError.Message = "Get to facebook access token url"
		oauthError.Error = err
		return nil, oauthError
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var userObj map[string]interface{}
	err = json.Unmarshal(body, &userObj)

	fmt.Println(userObj)

	b, err := json.MarshalIndent(userObj, "", "  ")
	logger.Info("Facebook User " + string(b))

	user.LastOauth = "FACEBOOK"
	email, _ := userObj["email"].(string)
	user.Email = email

	return user, nil
}
