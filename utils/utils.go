package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dmittelstaedt/dpaycol/models"
)

// GetExecDir sets the directory where the executable is located
func GetExecDir() string {
	exec, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(exec)
}

// GetServerID returns the id of the given name of the server
func GetServerID(apiURL string, hostname string) int {
	resp := SendRequest("GET", apiURL+"?"+hostname, nil)
	defer resp.Body.Close()

	var servers []models.Server
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		log.Println(err)
	}
	return servers[0].ID
}

// SendRequest sends request for given method and URL.
func SendRequest(method, url string, body []byte) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	return resp
}
