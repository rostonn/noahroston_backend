package main

import (
	"encoding/json"
	"go.uber.org/zap"
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
	configPath := strings.ToLower(os.Getenv("CONFIG_PATH"))
	if configPath == "" {
		log.Fatal("CONFIG_PATH env variable not set")
	}

	logger, err := createLogger(env)

	if err != nil {
		log.Fatal("Could not create Logger")
	}

	zap.ReplaceGlobals(logger)

	zap.S().Debug("Logger Created")

	configFilename := configPath + "config." + env + ".json"

	file, err := os.Open(configFilename)
	if err != nil {
		log.Fatal("Config Filename doesn't exist ", configFilename)
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
		zap.S().Error("Private Key ErRR")
		log.Fatal(sshErr)
	}

	sshPublicKeyBytes, e := ioutil.ReadFile(configuration.PublicKeyPath)
	if e != nil {
		zap.S().Error("Public Key ErRR")
		log.Fatal(sshErr)
	}

	zap.S().Info(string(sshPublicKeyBytes))

	verifyKey, err := jwt.ParseECPublicKeyFromPEM(sshPublicKeyBytes)

	if err != nil {
		zap.S().Error("Public Key Error")
		log.Fatal(err)
	}
	configuration.PublicKey = verifyKey
	configuration.PrivateKey = privateKey

	a := App{}
	a.Initialize(configuration)
	a.Run(":8080")
}
