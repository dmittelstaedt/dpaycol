package utils

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// GetExecDir sets the directory where the executable is located
func GetExecDir() string {
	exec, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(exec)
}

// SendRequest sends request for given method and URL.
func SendRequest(method, url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	return resp
}
