package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"go.uber.org/zap"
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
	PrivateKey           *ecdsa.PrivateKey
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func main() {
	rawJSON := []byte(`{
		"level": "debug",
		"development":true,
		"disableCaller":false,
		"encoding": "console",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "severity",
		  "levelEncoder": "uppercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	// "timeKey": "time",
	// "encodeTime": "zapcore.ISO8601TimeEncoder",
	// "callerKey":    "caller",
	// "encodeCaller": "zapcore.ShortCallerEncoder",
	cfg.EncoderConfig.TimeKey = "time"
	// cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// cfg.EncoderConfig.EncodeTime = TimeEncoder
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		s, ok := _levelToCapitalColorString[level]
		if !ok {
			s = _unknownLevelColor.Add(level.CapitalString())
		}
		enc.AppendString("[" + s + "]")
	}

	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	logger.Error("Logger error ...")

	slogger := logger.Sugar()
	slogger.Infow("Infow() allows tags", "name", "Legolas", "type", 1)

	flag.Parse()
	// flag.Lookup("log_dir").Value.Set("/Users/noahroston/go/src/github.com/rostonn/noahroston_backend/logs")
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Lookup("v").Value.Set("10")

	configuration := Configuration{}

	env := strings.ToLower(os.Getenv("PROFILE"))
	if env == "" {
		log.Fatal("PROFILE env variable not set")
	}
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
		fmt.Println("Private Key ErRR")
		log.Fatal(sshErr)
	}
	configuration.PrivateKey = privateKey

	a := App{}

	configJson, _ := json.Marshal(configuration)

	glog.V(2).Info("Confirguration:\n", string(configJson))
	// You need to set your Username and Password here
	a.Initialize(configuration)

	a.Run(":8080")
}
