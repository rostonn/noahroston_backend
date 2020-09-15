package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"

	"github.com/rostonn/noahroston_backend/models"
)

func getIpAddressInfo(accessKey, ipAddress string) models.UserLoginRecord {
	var ip models.UserLoginRecord

	var url string
	if ipAddress == "127.0.0.1" {
		zap.S().Debug("Local IP Address ")
		url = "http://api.ipstack.com/check?access_key=" + accessKey
	} else {
		zap.S().Info("Requestor IP Addres: " + ipAddress)
		url = "http://api.ipstack.com/" + ipAddress + "?access_key=" + accessKey
	}

	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error("IP FETCH ERROR ...")
		return ip
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &ip)
	return ip
}
