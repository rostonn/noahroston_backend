package main

import (
	"encoding/json"
	"fmt"
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
	logger.Debug("Logger Created")

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
		logger.Error("Private Key ErRR")
		log.Fatal(sshErr)
	}

	sshPublicKeyBytes, e := ioutil.ReadFile(configuration.PublicKeyPath)
	if e != nil {
		logger.Error("Public Key ErRR")
		log.Fatal(sshErr)
	}

	fmt.Println(string(sshPublicKeyBytes))

	// s := []byte("AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBIFCofI85aeUfgJ/qxY0f6aKCQNwCBA2GOyAk6y1+qqG4Rbw6OI67ZWvgLO7B/gvFtyZJLFThZHyaP38eO6qbc=")
	// verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(sshPublicKeyBytes)
	// 	s := []byte(`
	// -----BEGIN PUBLIC KEY-----
	// AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBIFCofI85aeUfgJ/qxY0f6aKCQNwCBA2GOyAk6y1+qqG4Rbw6OI67ZWvgLO7B/gvFtyZJLFThZHyaP38eO6qbc=
	// -----END PUBLIC KEY-----`)

	verifyKey, err := jwt.ParseECPublicKeyFromPEM(sshPublicKeyBytes)

	if err != nil {
		fmt.Println("Public Key Error")
		log.Fatal(err)
	}
	fmt.Println(verifyKey)
	// configuration.PublicKeyBytes = sshPublicKeyBytes
	configuration.PublicKey = verifyKey
	configuration.PrivateKey = privateKey

	a := App{}
	a.Initialize(configuration, logger)
	a.Run(":8080")
}
