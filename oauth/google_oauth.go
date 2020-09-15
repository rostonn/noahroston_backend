package oauth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rostonn/noahroston_backend/app"
	"github.com/rostonn/noahroston_backend/models"
	"go.uber.org/zap"
)

func LoginWithGoogle(code string, config app.Configuration) (*models.User, *OauthError) {
	zap.S().Debug("loginWithGoogle: google token", zap.String("code", code))
	user := &models.User{}
	oauthError := &OauthError{}

	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + code)

	if err != nil {
		zap.S().Error("Get to google token url")
		// panic(err)
		oauthError.Code = 500
		oauthError.Message = "Get to google access token url"
		oauthError.Error = err
		return nil, oauthError
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var userObj map[string]interface{}
	err = json.Unmarshal(body, &userObj)

	zap.S().Debug(userObj)

	b, err := json.MarshalIndent(userObj, "", "  ")
	zap.S().Info("Google User " + string(b))

	user.LastOauth = "GOOGLE"
	email, _ := userObj["email"].(string)
	user.Email = email

	return user, nil
}
