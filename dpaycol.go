package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"
)

// stats holds statistics of start or end of a job. IsEnd and RC are only used
// when the job finished.
type stats struct {
	Timestamp time.Time `json:"timestamp"`
	Ak        string    `json:"ak"`
	Am        int       `json:"am"`
	Lt        string    `json:"lt"`
	Ut        string    `json:"ut"`
	Jn        string    `json:"jn"`
	Jkid      string    `json:"jkid"`
	IsEnd     bool      `json:"isEnd"`
	RC        int       `json:"rc"`
}

const statsLog string = "stats-dpay.json"

var runStats stats

func init() {
	flag.StringVar(&runStats.Ak, "ak", "None", "Abrechnungskreis, required")
	flag.IntVar(&runStats.Am, "am", 0, "Abrechnungsmonat, required")
	flag.StringVar(&runStats.Lt, "lt", "None", "Lauftermin, optinal")
	flag.StringVar(&runStats.Ut, "ut", "None", "Untertermin, optional")
	flag.StringVar(&runStats.Jn, "jn", "None", "Jobname, required")
	flag.StringVar(&runStats.Jkid, "jkid", "None", "ID of the Jobkette, required")
	flag.BoolVar(&runStats.IsEnd, "e", false, "Ende, optional")
	flag.IntVar(&runStats.RC, "rc", -1, "Return Code of the Job, required if e")
}

// writeJSON writes the run statistics to a file in JSON encoding.
func (rs *stats) writeJSON() {
	bytes, err := json.Marshal(rs)
	if err != nil {
		log.Fatal(err)
	}

	bytes = append(bytes, "\n"...)

	file, err := os.OpenFile(statsLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := file.Write(bytes); err != nil {
		log.Fatal(err)
	}
}

// checkRequired checks if all required parameters are set
func (rs *stats) checkRequired() bool {
	if rs.Ak == "None" || rs.Am == 0 || rs.Jn == "None" || rs.Jkid == "None" {
		return false
	}
	return true
}

// checkMonth checks if month is between 1 and 12.
func (rs *stats) checkMonth() bool {
	if rs.Am < 1 || rs.Am > 12 {
		return false
	}
	return true
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

	if !runStats.checkRequired() {
		log.Fatal("Missing required parameters.")
	}

	if !runStats.checkMonth() {
		log.Fatal("Invalid month for am.")
	}

	if !runStats.checkEnd() {
		log.Fatal("Wrong parameters for end condition.")
	}

	runStats.writeJSON()
}
