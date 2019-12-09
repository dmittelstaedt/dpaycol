package models

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"regexp"
	"time"
)

// Stats holds statistics of start or end of a job. IsEnd and RC are only used
// when the job finished.
type Stats struct {
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

// writeJSON writes the run statistics to a file in JSON encoding.
func (rs *Stats) writeJSON(c Configuration, execDir string, statsFile string) {
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
func (rs *Stats) checkRequired() bool {
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
