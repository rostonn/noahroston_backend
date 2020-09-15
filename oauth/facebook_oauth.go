package oauth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rostonn/noahroston_backend/app"
	"github.com/rostonn/noahroston_backend/models"
	"go.uber.org/zap"
)

func LoginWithFacebook(code string, config app.Configuration) (*models.User, *OauthError) {
	zap.S().Debug("loginWithFacebook: token", zap.String("code", code))
	user := &models.User{}
	oauthError := &OauthError{}

	resp, err := http.Get("https://graph.facebook.com/v3.3/me?fields=email&access_token=" + code)

	if err != nil {
		zap.S().Error("Get to facebook token url")
		oauthError.Code = 500
		oauthError.Message = "Get to facebook access token url"
		oauthError.Error = err
		return nil, oauthError
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var userObj map[string]interface{}
	err = json.Unmarshal(body, &userObj)

	zap.S().Debug(userObj)

	b, err := json.MarshalIndent(userObj, "", "  ")
	zap.S().Info("Facebook User " + string(b))

	user.LastOauth = "FACEBOOK"
	email, _ := userObj["email"].(string)
	user.Email = email

	return user, nil
}
