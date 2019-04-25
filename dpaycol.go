package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"
)

const statsLog string = "stats-dpaycol.json"

// stats holds statistics of start or end of a job. IsEnd and RC are only used
// when the job finished.
type stats struct {
	Timestamp time.Time `json:"timestamp"`
	Ak        string    `json:"ak"`
	Am        string    `json:"am"`
	Lt        string    `json:"lt"`
	Ut        string    `json:"ut"`
	Jn        string    `json:"jn"`
	Jkid      string    `json:"jkid"`
	IsEnd     bool      `json:"isEnd"`
	RC        int       `json:"rc"`
}

var runStats stats

func init() {
	flag.StringVar(&runStats.Ak, "ak", "None", "Abrechnungskreis, required")
	flag.StringVar(&runStats.Am, "am", "None", "Abrechnungsmonat, required")
	flag.StringVar(&runStats.Lt, "lt", "None", "Lauftermin, required")
	flag.StringVar(&runStats.Ut, "ut", "None", "Untertermin, required")
	flag.StringVar(&runStats.Jn, "jn", "None", "Jobname, required")
	flag.StringVar(&runStats.Jkid, "jkid", "None", "ID of the Jobkette, required")
	flag.BoolVar(&runStats.IsEnd, "e", false, "Ende, optional")
	flag.IntVar(&runStats.RC, "rc", -1, "Return Code of the Job, required if e")
	flag.Parse()
	runStats.Timestamp = time.Now()
}

// writeJSON writes the run statistics to a file in JSON encoding.
func writeJSON() {
	bytes, err := json.Marshal(runStats)
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
func checkRequired() bool {
	if runStats.Ak == "None" || runStats.Am == "None" || runStats.Jn == "None" || runStats.Jkid == "None" {
		return false
	}
	return true
}

// checkEnd checks if end condition is correct set. Correct means either no e
// and rc is set or e is set and rc is >= 0.
func checkEnd() bool {
	if runStats.IsEnd == false && runStats.RC < 0 {
		return true
	}
	if runStats.IsEnd == true && runStats.RC >= 0 {
		return true
	}
	return false
}

func main() {
	if required := checkRequired(); required == false {
		log.Fatal("Missing required parameters.")
	}

	if end := checkEnd(); end == false {
		log.Fatal("Wrong parameters for end condition.")
	}

	writeJSON()
}
