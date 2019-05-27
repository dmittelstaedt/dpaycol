package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/spf13/viper"
)

// stats holds statistics of start or end of a job. IsEnd and RC are only used
// when the job finished.
type stats struct {
	Timestamp time.Time `json:"timestamp"`
	HostName  string    `json:"hostname"`
	Ak        string    `json:"ak"`
	Am        string    `json:"am"`
	Lt        string    `json:"lt"`
	Ut        string    `json:"ut"`
	Jn        string    `json:"jn"`
	Jkid      string    `json:"jkid"`
	IsEnd     bool      `json:"isEnd"`
	RC        int       `json:"rc"`
}

type configuration struct {
	StatsPath string
}

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

var runStats stats
var version *bool
var versionNumber string
var gitCommit string
var buildDate string
var execDir string

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
func setExecDir() {
	exec, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDir = filepath.Dir(exec)
}

// readConfig reads configuration from file
func readConfig() configuration {
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

// writeJSON writes the run statistics to a file in JSON encoding.
func (rs *stats) writeJSON(c configuration) {
	bytes, err := json.Marshal(rs)
	if err != nil {
		log.Fatal(err)
	}

	bytes = append(bytes, "\n"...)

	var statsLogPath string
	if c.StatsPath == "" {
		statsLogPath = execDir + "/" + statsFile
	} else {
		statsLogPath = c.StatsPath
	}

	file, err := os.OpenFile(statsLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		log.Fatal(err)
	}
}

// checkRequired checks if all required parameters are set
func (rs *stats) checkRequired() bool {
	requiredAk := false
	requiredAm := false
	requiredJn := false
	requiredJkid := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == ak {
			requiredAk = true
		}
		if f.Name == am {
			requiredAm = true
		}
		if f.Name == jn {
			requiredJn = true
		}
		if f.Name == jkid {
			requiredJkid = true
		}
	})
	return requiredAk && requiredAm && requiredJn && requiredJkid
}

// checkMonth checks if month is between 1 and 12.
func (rs *stats) checkMonth() bool {
	r, err := regexp.Compile("^\\d{4}(0[1-9]|1[0-2])$")
	if err != nil {
		return false
	}
	return r.MatchString(rs.Am)
}

// checkEnd checks if end condition is correct set. Correct means either no e
// and rc is set or e is set and rc is >= 0.
func (rs *stats) checkEnd() bool {
	if rs.IsEnd == false && rs.RC < 0 {
		return true
	}
	if rs.IsEnd == true && rs.RC >= 0 {
		return true
	}
	return false
}

func main() {
	flag.Parse()
	runStats.Timestamp = time.Now()
	runStats.HostName, _ = os.Hostname()

	setExecDir()

	configuration := readConfig()

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
