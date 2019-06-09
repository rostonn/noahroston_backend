package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

var a = App{}

var logger = zap.NewExample()
var configuration = Configuration{}

func TestMain(t *testing.T) {
	a.Initialize(configuration, logger)

	req, _ := http.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	compare(t, 200, rr.Code)
	compare(t, rr.Body.String(), "OK")
}

func compare(t *testing.T, expected, received interface{}) {
	if expected != received {
		t.Errorf("Expected=%v (%T) Received=%v (%T)", expected, expected, received, received)
	}
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}
func clearTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS users
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age INT NOT NULL
)`

// configFilename := "config/config.test.json"

// file, err := os.Open(configFilename)
// if err != nil {
// 	log.Fatal("Test Config Filename doesn't exist", configFilename)
// }
// decoder := json.NewDecoder(file)
// err = decoder.Decode(&configuration)
// if err != nil {
// 	log.Fatal(err)
// }

// sshPrivateKeyBytes, sshErr := ioutil.ReadFile(configuration.PrivateKeyPath)
// if sshErr != nil {
// 	log.Fatal(sshErr)
// }

// privateKey, sshErr := jwt.ParseECPrivateKeyFromPEM(sshPrivateKeyBytes)
// if sshErr != nil {
// 	logger.Error("Private Key ErRR")
// 	log.Fatal(sshErr)
// }
// configuration.PrivateKey = privateKey
