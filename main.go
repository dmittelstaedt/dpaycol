package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dmittelstaedt/dpaycol/models"

	"github.com/spf13/viper"
)

// constantes for log file and name of the parameters
const (
	configFile string = "config"
	statsFile  string = "stats-dpay.json"
	ak         string = "ak"
	am         string = "am"
	ut         string = "ut"
	lt         string = "lt"
	jn         string = "jn"
	jkid       string = "jkid"
	e          string = "e"
	rc         string = "rc"
	v          string = "v"
)

var runStats models.Stats
var version *bool
var versionNumber string
var gitCommit string
var buildDate string

func init() {
	flag.StringVar(&runStats.Ak, ak, "", "Abrechnungskreis, required")
	flag.StringVar(&runStats.Am, am, "", "Abrechnungsmonat (YYYYMM), required")
	flag.StringVar(&runStats.Lt, lt, "None", "Lauftermin, optinal")
	flag.StringVar(&runStats.Ut, ut, "None", "Untertermin, optional")
	flag.StringVar(&runStats.Jn, jn, "", "Jobname, required")
	flag.StringVar(&runStats.Jkid, jkid, "", "ID of the Jobkette, required")
	flag.BoolVar(&runStats.IsEnd, e, false, "Ende, optional")
	flag.IntVar(&runStats.RC, rc, -1, "Return Code of the Job, required if e")
	version = flag.Bool(v, false, "Version")
}

// execDir sets the directory where the executable is located
func setExecDir() string {
	exec, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(exec)
}

// readConfig reads configuration from file
func readConfig() models.Configuration {
	viper.SetConfigName(configFile)
	viper.AddConfigPath(execDir)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		} else {
			return configuration{
				StatsPath: "",
			}
		}
	}

	var configuration configuration
	err = viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	return configuration
}

func checkRestAPI(configuration configuration) {
	apiURL := "http://" + configuration.APIEndpoint + "/servers"
	log.Printf("API URL: %+v", apiURL)
	resp := sendRequest("GET", apiURL)
	defer resp.Body.Close()

	var servers []models.Server
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		log.Println(err)
	}

	log.Printf("Server: %+v", servers)
}

// sendRequest sends request for given method and URL.
func sendRequest(method, url string) *http.Response {
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

func main() {
	flag.Parse()
	runStats.Timestamp = time.Now()
	runStats.HostName, _ = os.Hostname()

	execDir := setExecDir()

	configuration := readConfig()

	checkRestAPI(configuration)

	if *version {
		fmt.Println("Version: " + versionNumber)
		fmt.Println("Git Commit: " + gitCommit)
		fmt.Println("Build Date: " + buildDate)
		os.Exit(0)
	}

	if !runStats.checkRequired() {
		fmt.Println("Missing required parameters. Use -h for usage.")
		os.Exit(1)
	}

	if !runStats.checkMonth() {
		fmt.Println("Invalid month for am. Use -h for usage.")
		os.Exit(2)
	}

	if !runStats.checkEnd() {
		fmt.Println("Wrong parameters for end condition. Use -h for usage.")
		os.Exit(3)
	}

	runStats.writeJSON(configuration)
}
