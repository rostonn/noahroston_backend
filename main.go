package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/rostonn/noahroston_backend/app"
)

// type Configuration struct {
// 	DbUsername           string `json:"db_username"`
// 	DbPassword           string `json:"db_password"`
// 	Dbname               string `json:"db_name"`
// 	Port                 string `json:"port"`
// 	AmazonClientId       string `json:"amazonClientId"`
// 	AmazonClientSecret   string `json:"amazonClientSecret"`
// 	AmazonAccessTokenURL string `json:"amazonAccessTokenURL"`
// 	AmazonProfileURL     string `json:"amazonProfileURL"`
// 	AmazonRedirectURI    string `json:"amazonRedirectURI"`
// 	IpStackApiKey        string `json:"ipStackApiKey"`
// 	PrivateKeyPath       string `json:"privateKeyPath"`
// 	PrivateKey           *ecdsa.PrivateKey
// }

func main() {

	configuration := app.Configuration{}

	env := strings.ToLower(os.Getenv("PROFILE"))
	if env == "" {
		log.Fatal("PROFILE env variable not set")
	}

	logger, err := createLogger(env)
	if err != nil {
		log.Fatal("Could not create Logger")
	}
	logger.Debug("Logger Created")

	configFilename := "config/config." + env + ".json"

	file, err := os.Open(configFilename)
	if err != nil {
		log.Fatal("Config Filename doesn't exist", configFilename)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	// read id_rsa file
	sshPrivateKeyBytes, sshErr := ioutil.ReadFile(configuration.PrivateKeyPath)
	if sshErr != nil {
		log.Fatal(sshErr)
	}
	privateKey, sshErr := jwt.ParseECPrivateKeyFromPEM(sshPrivateKeyBytes)
	if sshErr != nil {
		logger.Error("Private Key ErRR")
		log.Fatal(sshErr)
	}
	configuration.PrivateKey = privateKey

	a := App{}
	a.Initialize(configuration, logger)
	a.Run(":8080")
}
