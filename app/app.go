package app

import (
	"crypto/ecdsa"
)

type Configuration struct {
	DbUsername           string `json:"db_username"`
	DbPassword           string `json:"db_password"`
	Dbname               string `json:"db_name"`
	Port                 string `json:"port"`
	AmazonClientId       string `json:"amazonClientId"`
	AmazonClientSecret   string `json:"amazonClientSecret"`
	AmazonAccessTokenURL string `json:"amazonAccessTokenURL"`
	AmazonProfileURL     string `json:"amazonProfileURL"`
	AmazonRedirectURI    string `json:"amazonRedirectURI"`
	IpStackApiKey        string `json:"ipStackApiKey"`
	PrivateKeyPath       string `json:"privateKeyPath"`
	PublicKeyPath        string `json:"publicKeyPath"`
	PrivateKey           *ecdsa.PrivateKey
	PublicKey            *ecdsa.PublicKey
}
